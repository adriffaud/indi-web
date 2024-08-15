defmodule IndiEx.Repo do
  use Ecto.Repo,
    otp_app: :indi_ex,
    adapter: Ecto.Adapters.SQLite3
end
