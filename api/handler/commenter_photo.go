package handler

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"simple-commenting/repository"
	"strings"

	"github.com/disintegration/imaging"
)

func CommenterPhotoHandler(w http.ResponseWriter, r *http.Request) {
	commenter, err := repository.Repo.CommenterRepository.GetCommenterByHex(r.FormValue("commenterHex"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	url := commenter.Photo
	if commenter.Provider == "google" {
		if strings.HasSuffix(url, "photo.jpg") {
			url += "?sz=38"
		} else if strings.Contains(url, "=") {
			url = strings.Split(url, "=")[0] + "=s38"
		} else {
			url += "=s38"
		}
	} else if commenter.Provider == "github" {
		url += "&s=38"
	} else if commenter.Provider == "gitlab" {
		url += "?width=38"
	}

	resp, err := http.Get(url)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer resp.Body.Close()

	if commenter.Provider != "commento" { // Custom URL avatars need to be resized.
		io.Copy(w, resp.Body)
		return
	}

	// Limit the size of the response to 128 KiB to prevent DoS attacks
	// that exhaust memory.
	limitedResp := &io.LimitedReader{R: resp.Body, N: 128 * 1024}

	img, _, err := image.Decode(limitedResp)
	if err != nil {
		fmt.Fprintf(w, "Image decode failed: %v\n", err)
		return
	}

	if err = imaging.Encode(w, imaging.Resize(img, 38, 0, imaging.Lanczos), imaging.JPEG); err != nil {
		fmt.Fprintf(w, "image encoding failed: %v\n", err)
		return
	}
}
