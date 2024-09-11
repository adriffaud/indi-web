defmodule IndiExWeb.PageController do
  use IndiExWeb, :controller

  def home(conn, _params) do
    render(conn, :home)
  end
end
