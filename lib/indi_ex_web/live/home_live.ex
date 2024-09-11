defmodule IndiExWeb.HomeLive do
  use IndiExWeb, :live_view

  require Logger

  @impl true
  def mount(_params, _session, socket) do
    {:ok, socket}
  end

  @impl true
  def handle_event(event, _unsigned_params, socket) do
    Logger.debug("Unhandled event: #{inspect(event)}")
    {:noreply, socket}
  end
end
