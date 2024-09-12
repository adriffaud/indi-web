defmodule IndiEx.IndiClient do
  use GenServer

  require Logger

  @initial_state %{socket: nil, timer: nil, partial: nil, properties: []}

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

    new_timer = Process.send_after(self(), {:tcp_idle, socket}, 1000)
    {:ok, partial} = Saxy.Partial.new(IndiEx.IndiXml, [])
    {:ok, %{state | socket: socket, timer: new_timer, partial: partial}}
  end

  @impl true
  def handle_cast({:command, cmd}, %{socket: socket} = state) do
    :ok = :gen_tcp.send(socket, cmd)
    {:noreply, state}
  end

  @impl true
  def handle_info({:tcp, socket, msg}, %{partial: partial} = state) do
    Process.cancel_timer(state.timer)

    dbg(msg)
    res = Saxy.Partial.parse(partial, msg)
    dbg(res)
    {:cont, partial} = res

    new_timer = Process.send_after(self(), {:tcp_idle, socket}, 1000)
    {:noreply, %{state | partial: partial, timer: new_timer}}
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
