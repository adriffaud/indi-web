defmodule IndiExWeb.HardwareLive do
  use IndiExWeb, :live_view

  require Logger

  @properties [
    %{
      device: "Telescope Simulator",
      group: "Connection",
      type: :switch,
      name: "CONNECTION",
      label: "Connection",
      perm: "rw",
      values: [
        %{name: "CONNECT", label: "Connect", value: "On"},
        %{name: "DISCONNECT", label: "Disconnect", value: "Off"}
      ]
    },
    %{
      device: "Telescope Simulator",
      group: "Main control",
      type: :text,
      name: "CONNECTION",
      label: "Connection",
      perm: "rw",
      values: [
        %{name: "CONNECT", label: "Connect", value: "On"},
        %{name: "DISCONNECT", label: "Disconnect", value: "Off"}
      ]
    }
  ]

  @impl true
  def mount(_params, _session, socket) do
    tree =
      @properties
      |> Enum.group_by(& &1.device)
      |> Enum.map(fn {device, props} ->
        {device, Enum.group_by(props, & &1.group)}
      end)
      |> Enum.into(%{})
      |> IO.inspect()

    socket = assign(socket, :properties, tree)
    {:ok, socket}
  end
end
