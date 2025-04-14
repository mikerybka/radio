package radio

import (
	"net/http"
	"os"
	"os/exec"
)

func NewServer() http.Handler {
	radio := http.NewServeMux()
	radio.HandleFunc("GET /live.mp3", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/mpeg")
		cmd := exec.Command("ffmpeg",
			"-f", "pulse",
			"-i", "virtual_output.monitor",
			"-f", "mp3",
			"pipe:1",
		)
		cmd.Stdout = w
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	radio.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<audio controls autoplay>
			<source src="/live.mp3" type="audio/mpeg">
			Your browser does not support the audio element.
		</audio>`))
	})
	return radio
}
