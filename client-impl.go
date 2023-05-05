package openai_scheduler

import (
	"context"
	o "github.com/sashabaranov/go-openai"
)

// region ChatCompletion

func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request o.ChatCompletionRequest,
) (
	response o.ChatCompletionResponse, err error,
) {
	response, err = c.Raw.CreateChatCompletion(ctx, request)
	c.StatusAdjust(err)
	return
}

func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request o.ChatCompletionRequest,
) (
	stream *o.ChatCompletionStream, err error,
) {
	stream, err = c.Raw.CreateChatCompletionStream(ctx, request)
	c.StatusAdjust(err)
	return
}

// endregion

// region Completion

func (c *Client) CreateCompletion(
	ctx context.Context,
	request o.CompletionRequest,
) (
	response o.CompletionResponse, err error,
) {
	response, err = c.Raw.CreateCompletion(ctx, request)
	c.StatusAdjust(err)
	return
}

func (c *Client) CreateCompletionStream(
	ctx context.Context, request o.CompletionRequest,
) (
	stream *o.CompletionStream, err error,
) {
	stream, err = c.Raw.CreateCompletionStream(ctx, request)
	c.StatusAdjust(err)
	return
}

// endregion

// region FineTune

func (c *Client) CreateFineTune(ctx context.Context, request o.FineTuneRequest) (response o.FineTune, err error) {
	response, err = c.Raw.CreateFineTune(ctx, request)
	c.StatusAdjust(err)
	return
}

func (c *Client) CancelFineTune(ctx context.Context, fineTuneID string) (response o.FineTune, err error) {
	response, err = c.Raw.CancelFineTune(ctx, fineTuneID)
	c.StatusAdjust(err)
	return
}

func (c *Client) ListFineTunes(ctx context.Context) (response o.FineTuneList, err error) {
	response, err = c.Raw.ListFineTunes(ctx)
	c.StatusAdjust(err)
	return
}

func (c *Client) GetFineTune(ctx context.Context, fineTuneID string) (response o.FineTune, err error) {
	response, err = c.Raw.GetFineTune(ctx, fineTuneID)
	c.StatusAdjust(err)
	return
}

func (c *Client) DeleteFineTune(ctx context.Context, fineTuneID string) (response o.FineTuneDeleteResponse, err error) {
	response, err = c.Raw.DeleteFineTune(ctx, fineTuneID)
	c.StatusAdjust(err)
	return
}

func (c *Client) ListFineTuneEvents(ctx context.Context, fineTuneID string) (response o.FineTuneEventList, err error) {
	response, err = c.Raw.ListFineTuneEvents(ctx, fineTuneID)
	c.StatusAdjust(err)
	return
}

// endregion

// region Image

func (c *Client) CreateImage(ctx context.Context, request o.ImageRequest) (response o.ImageResponse, err error) {
	response, err = c.Raw.CreateImage(ctx, request)
	c.StatusAdjust(err)
	return
}

func (c *Client) CreateEditImage(ctx context.Context, request o.ImageEditRequest) (response o.ImageResponse, err error) {
	response, err = c.Raw.CreateEditImage(ctx, request)
	c.StatusAdjust(err)
	return
}

func (c *Client) CreateVariImage(ctx context.Context, request o.ImageVariRequest) (response o.ImageResponse, err error) {
	response, err = c.Raw.CreateVariImage(ctx, request)
	c.StatusAdjust(err)
	return
}

// endregion

// region Audio

func (c *Client) CreateTranscription(
	ctx context.Context,
	request o.AudioRequest,
) (response o.AudioResponse, err error) {
	response, err = c.Raw.CreateTranscription(ctx, request)
	c.StatusAdjust(err)
	return
}

func (c *Client) CreateTranslation(
	ctx context.Context,
	request o.AudioRequest,
) (response o.AudioResponse, err error) {
	response, err = c.Raw.CreateTranslation(ctx, request)
	c.StatusAdjust(err)
	return
}

// endregion

// region File

func (c *Client) CreateFile(ctx context.Context, request o.FileRequest) (file o.File, err error) {
	file, err = c.Raw.CreateFile(ctx, request)
	c.StatusAdjust(err)
	return
}

func (c *Client) DeleteFile(ctx context.Context, fileID string) (err error) {
	err = c.Raw.DeleteFile(ctx, fileID)
	c.StatusAdjust(err)
	return
}

func (c *Client) ListFiles(ctx context.Context) (files o.FilesList, err error) {
	files, err = c.Raw.ListFiles(ctx)
	c.StatusAdjust(err)
	return
}

func (c *Client) GetFile(ctx context.Context, fileID string) (file o.File, err error) {
	file, err = c.Raw.GetFile(ctx, fileID)
	c.StatusAdjust(err)
	return
}

// endregion

// region Engine

func (c *Client) ListEngines(ctx context.Context) (engines o.EnginesList, err error) {
	engines, err = c.Raw.ListEngines(ctx)
	c.StatusAdjust(err)
	return
}

func (c *Client) GetEngine(
	ctx context.Context,
	engineID string,
) (engine o.Engine, err error) {
	engine, err = c.Raw.GetEngine(ctx, engineID)
	c.StatusAdjust(err)
	return
}

// endregion

// region Moderations

func (c *Client) Moderations(ctx context.Context, request o.ModerationRequest) (response o.ModerationResponse, err error) {
	response, err = c.Raw.Moderations(ctx, request)
	c.StatusAdjust(err)
	return
}

// endregion
