package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/iDominate/golang-pexels-api/env"
)

const PHOTOS_URL_KEY = "PEXELS_PHOTOS_URL"
const VIDEOS_URL_KEY = "PEXELS_VIDEOS_URL"

type Client struct {
	Token           string
	HttpClient      http.Client
	ReamainingTimes int64
}

func (c *Client) NewClient(token string) *Client {
	return &Client{Token: token, HttpClient: http.Client{}}
}

func (c *Client) SearchForPhotos(query string, perPage int32, page int32) (*SearchResult, error) {
	photos_url, err := env.Getenv(PHOTOS_URL_KEY)
	if err != nil {
		log.Fatal(err.Error())
	}

	url := fmt.Sprintf("%s/search?query=%s&per_page=%d&page=%d", photos_url, query, perPage, page)

	response, err := c.NewRequest(http.MethodGet, url)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err.Error())
	}

	var result SearchResult

	err = json.Unmarshal(data, &result)

	if err != nil {
		log.Fatal(err)
	}

	return &result, err
}

func (c *Client) CuratedPhotos(perPage int32, page int32) (*CuratedResult, error) {
	key, err := env.Getenv(PHOTOS_URL_KEY)

	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/curated?per_page=%d&page=%d", key, perPage, page)

	response, err := c.NewRequest(http.MethodGet, url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	var result CuratedResult

	err = json.Unmarshal(data, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetPhotoById(id int32) (*Photo, error) {
	key, err := env.Getenv(PHOTOS_URL_KEY)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/photos/%d", key, id)
	res, err := c.NewRequest(http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var photo Photo
	err = json.Unmarshal(data, &photo)

	if err != nil {
		return nil, err
	}

	return &photo, nil
}

func (c *Client) NewRequest(method string, url string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", c.Token)

	response, err := c.HttpClient.Do(request)

	if err != nil {
		return response, err
	}
	times, err := strconv.Atoi(response.Header.Get("X-Ratelimit-Remaining"))

	if err != nil {
		return response, err
	}
	c.ReamainingTimes = int64(times)
	return response, nil
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	src := rand.New(rand.NewSource(time.Now().Unix()))

	id := src.Intn(10000000)

	response, err := c.GetPhotoById(int32(id))

	if err != nil {
		return response, err
	}
	return response, nil
}

func (c *Client) GetRemainingRequests() int32 {
	return int32(c.ReamainingTimes)
}

func (c *Client) SearchForVideo(query string, page int32, perPage int32) (*VideoSearchResult, error) {
	key, err := env.Getenv(VIDEOS_URL_KEY)

	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/search?query=%s&page=%d&per_page=%d", key, query, page, perPage)

	res, err := c.NewRequest(http.MethodGet, url)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	var VideoSearchResult VideoSearchResult

	err = json.Unmarshal(data, &VideoSearchResult)

	if err != nil {
		return nil, err
	}

	return &VideoSearchResult, nil
}

func (c *Client) GetPopularVideos(page, perPage int32) (*PopularVideos, error) {
	key, err := env.Getenv(VIDEOS_URL_KEY)

	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/popular?per_page=%d&page=%d", key, perPage, page)

	res, err := c.NewRequest(http.MethodGet, url)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	var popular PopularVideos
	err = json.Unmarshal(data, &popular)

	if err != nil {
		return nil, err
	}

	return &popular, nil
}

func (c *Client) GetRandomVideo() (*Video, error) {
	src := rand.New(rand.NewSource(time.Now().Unix()))
	random := src.Intn(10000000)

	res, err := c.GetPopularVideos(1, int32(random))

	if err != nil {
		return nil, err
	}
	return &res.Videos[0], nil
}
