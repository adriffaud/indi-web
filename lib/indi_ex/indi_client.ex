defmodule IndiEx.IndiClient do
  use GenServer

  require Logger

  @initial_state %{socket: nil, stream: nil, properties: []}

  def start_link do
    GenServer.start_link(__MODULE__, @initial_state)
  end

  def get_properties(pid) do
    GenServer.cast(pid, {:command, "<getProperties version=\"1.7\" />"})
  end

  @impl true
  def init(state) do
    opts = [:binary, active: true]
    {:ok, socket} = :gen_tcp.connect(~c"localhost", 7624, opts)
    {:ok, %{state | socket: socket}}
  end

  @impl true
  def handle_cast({:command, cmd}, %{socket: socket} = state) do
    :ok = :gen_tcp.send(socket, cmd)

    {:noreply, state}
  end

  @impl true
  def handle_info({:tcp, _socket, msg}, state) do
    Logger.debug("📬📬 Received data")
    IO.inspect(msg, pretty: true)

    # msg
    # |> Saxy.parse_string(IndiEx.IndiXml, [])
    # |> IO.inspect(pretty: true)

    {:noreply, state}
  end

  @impl true
  def handle_info(msg, state) do
    Logger.debug("⚠️⚠️⚠️ Unhandled message: #{inspect(msg)}")
    {:noreply, state}
  end
end
