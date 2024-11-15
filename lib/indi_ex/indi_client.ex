defmodule IndiEx.IndiClient do
  use GenServer, restart: :transient

  require Logger

  @initial_state %{socket: nil, timer: nil, partial: nil, properties: []}
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
    {:ok, socket} = :gen_tcp.connect(~c"localhost", 7624, opts)

    new_timer = Process.send_after(self(), {:tcp_idle, socket}, @idle_timeout)
    {:ok, partial} = Saxy.Partial.new(IndiEx.IndiXml, [])
    {:ok, %{state | socket: socket, timer: new_timer, partial: partial}}
  end

  @impl true
  def handle_cast({:command, cmd}, %{socket: socket} = state) do
    Logger.debug("📬📬 Sending command: #{cmd}")
    :ok = :gen_tcp.send(socket, cmd)
    {:noreply, state}
  end

  @impl true
  def handle_info({:tcp, socket, msg}, %{partial: partial} = state) do
    Process.cancel_timer(state.timer)

    new_timer = Process.send_after(self(), {:tcp_idle, socket}, @idle_timeout)

    case Saxy.Partial.parse(partial, msg) do
      {:cont, partial} ->
        {:noreply, %{state | partial: partial, timer: new_timer}}
      {:error, err} ->
        dbg(err)
        {:noreply, %{state |  timer: new_timer}}
    end

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
end
