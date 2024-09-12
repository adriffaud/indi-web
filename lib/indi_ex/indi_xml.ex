defmodule IndiEx.IndiXml do
  @behaviour Saxy.Handler

  def handle_event(:start_document, prolog, state) do
    {:ok, [{:start_document, prolog} | state]}
  end

  def handle_event(:end_document, _data, state) do
    {:ok, [{:end_document} | state]}
  end

  def handle_event(:start_element, {name, attributes}, state) do
    {:ok, [{:start_element, name, attributes} | state]}
  end

  def handle_event(:end_element, name, state) do
    {:ok, [{:end_element, name} | state]}
  end

  def handle_event(:characters, chars, state) do
    chars = String.trim(chars)
    {:ok, [{:characters, chars} | state]}
  end

  def handle_event(:cdata, cdata, state) do
    cdata = String.trim(cdata)
    {:ok, [{:cdata, cdata} | state]}
  end
end
