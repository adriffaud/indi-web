package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	"github.com/adriffaud/indi-web/internal/mount"
	"github.com/adriffaud/indi-web/ui/components"
	"github.com/adriffaud/indi-web/ui/pages"
)

type SSEClient struct {
	eventChan  chan templ.Component
	indiClient *indiclient.Client
	mount      mount.Mount
}

func (sse SSEClient) OnNotify(e indiclient.Event) {
	var renderList []templ.Component

	switch e.EventType {
	case indiclient.Timeout:
		return
	case indiclient.Add, indiclient.Delete:
		renderList = append(renderList, pages.DeviceView(sse.indiClient.Properties, e.Property.Device))
	case indiclient.Update:
		if e.Property.Device == sse.mount.Driver && e.Property.Name == "EQUATORIAL_EOD_COORD" {
			for _, value := range e.Property.Values {
				if value.Name == "RA" {
					formated, err := DecimalToSexagesimal(value.Value)
					if err != nil {
						slog.Error("could not convert RA coords in sexagesimal", "err", err, "value", value)
					}
					sse.mount.RA = formated
					renderList = append(renderList, components.TextInput(
						"ra_input",
						sse.mount.RA,
						templ.Attributes{"disabled": "true", "hx-swap-oob": "true"},
					))
					if err != nil {
						slog.Error("failed to convert to HTML", "error", err)
					}
				}
				if value.Name == "DEC" {
					formated, err := DecimalToSexagesimal(value.Value)
					if err != nil {
						slog.Error("could not convert DEC coords in sexagesimal", "err", err, "value", value)
					}
					sse.mount.DEC = formated
					renderList = append(renderList, components.TextInput(
						"dec_input",
						sse.mount.DEC,
						templ.Attributes{"disabled": "true", "hx-swap-oob": "true"},
					))
					if err != nil {
						slog.Error("failed to convert to HTML", "error", err)
					}
				}
			}
		}

		renderList = append(renderList, pages.PropertyValues(e.Property))
	case indiclient.Message:
		slog.Debug("📮 Notification", "message", e.Message)
		renderList = append(renderList, pages.Notifications(e.Message))
	}

	if renderList == nil {
		return
	}

	for _, comp := range renderList {
		sse.eventChan <- comp
	}
}

func (app *application) sse(w http.ResponseWriter, r *http.Request) {
	slog.Debug("✅ SSE client connected", "address", r.RemoteAddr, "referer", r.Referer())

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	client := SSEClient{eventChan: make(chan templ.Component), indiClient: app.indiClient, mount: app.mount}
	app.indiClient.Register(client)

	for {
		select {
		case <-r.Context().Done():
			slog.Debug("SSE client disconnected", "address", r.RemoteAddr)
			app.indiClient.Unregister(client)
			return
		case component := <-client.eventChan:
			var buf bytes.Buffer
			var err error

			component.Render(r.Context(), &buf)
			if err != nil {
				slog.Error("failed to convert to HTML", "error", err)
			}

			fmt.Fprintf(w, "data: %s\n\n", buf.String())
			w.(http.Flusher).Flush()
		}
	}
}

func DecimalToSexagesimal(decimal string) (string, error) {
	decimalf, err := strconv.ParseFloat(decimal, 64)
	if err != nil {
		return "", err
	}

	hours, remainder := math.Modf(decimalf)
	minutesf := remainder * 60
	minutes, remainder := math.Modf(minutesf)
	seconds := int(remainder * 60)

	return fmt.Sprintf("%02.0f:%02.0f:%02d", hours, minutes, seconds), nil
}
