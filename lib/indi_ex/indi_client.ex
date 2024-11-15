defmodule IndiEx.IndiClient do
  use GenServer, restart: :transient

  require Logger

  @initial_state %{socket: nil, timer: nil, buffer: "", properties: []}
  @idle_timeout 5000

  def start_link do
    GenServer.start_link(__MODULE__, @initial_state, name: __MODULE__)
  end

  def get_properties() do
    GenServer.cast(__MODULE__, {:command, "<getProperties version=\"1.7\" />"})
  end

  @impl true
  def init(state) do
    opts = [:binary, active: true]

    case :gen_tcp.connect(~c"localhost", 7624, opts) do
      {:ok, socket} ->
        new_timer = Process.send_after(self(), {:tcp_idle, socket}, @idle_timeout)
        {:ok, %{state | socket: socket, timer: new_timer}}

      {:error, reason} ->
        Logger.error("Failed to connect to indiserver - #{reason}")
        {:stop, reason}
    end
  end

  @impl true
  def handle_cast({:command, cmd}, %{socket: socket} = state) do
    Logger.debug("📬📬 Sending command: #{cmd}")
    :ok = :gen_tcp.send(socket, cmd)
    {:noreply, state}
  end

  @impl true
  def handle_info({:tcp, socket, data}, %{buffer: buffer} = state) do
    Process.cancel_timer(state.timer)

    new_buffer = buffer <> data
    {complete_elements, leftover} = process_buffer(new_buffer)

    Enum.each(complete_elements, fn element ->
      dbg(element)
    end)

    new_timer = Process.send_after(self(), {:tcp_idle, socket}, @idle_timeout)
    {:noreply, %{state | timer: new_timer, buffer: leftover}}
  end

  @impl true
  def handle_info({:tcp_idle, _socket}, state) do
    Logger.debug("😴😴 Connection idle")
    {:noreply, state}
  end

  @impl true
  def handle_info(msg, state) do
    Logger.debug("⚠️⚠️⚠️ Unhandled message: #{inspect(msg)}")
    {:noreply, state}
  end

  defp process_buffer(buffer) do
    case Saxy.parse_string(buffer, IndiEx.IndiXml, []) do
      {:ok, data} ->
        dbg(data)

      {:error, %Saxy.ParseError{} = error} ->
        dbg(error)
    end

    {[], buffer}
  end
end
