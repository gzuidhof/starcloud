/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package npmcdn

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starcloud/starcloud/logger"
	"github.com/gzuidhof/starcloud/starcloud/npm"
)

type NPMCDNHandler struct {
	root http.FileSystem
	prefixToStrip string
	cacheControlHeaderValue string
}

func NewNPMCDNHandler (prefix string, root http.FileSystem) *NPMCDNHandler {
	return &NPMCDNHandler{
		prefixToStrip: prefix,
		root: root,
		cacheControlHeaderValue: "public, max-age=604800",
	}
}

// Based on https://github.com/gofiber/fiber/blob/master/middleware/filesystem/filesystem.go
func (h *NPMCDNHandler) Handler(c *fiber.Ctx) error {
	path := strings.TrimPrefix(c.Path(), h.prefixToStrip)

	pkg, err := npm.PathToNPMPackage(path)

	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Not found: %s", err.Error()))
	}

	if (pkg.Name != "starboard-notebook") {
		return c.Status(fiber.StatusNotFound).SendString("Only starboard-notebook is currently served")
	}

	file, err := h.root.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		logger.GetSugarLogger(c).Warnf("Failed to open path %s", path)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
	}

	stat, err := file.Stat()
	if err != nil {
		logger.GetSugarLogger(c).Warnf("Failed to stat path %s", path)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to stat file")
	}

	if stat.IsDir() {
		return dirList(c, file)
	}

	modTime := stat.ModTime()
	contentLength := int(stat.Size())

	c.Type(getFileExtension(stat.Name()))

	// Set Last Modified header
	if !modTime.IsZero() {
		c.Set(fiber.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
	}

	method := c.Method()
	if method == fiber.MethodGet {
		c.Set(fiber.HeaderCacheControl, h.cacheControlHeaderValue)
		c.Response().SetBodyStream(file, contentLength)
		return nil
	}
	if method == fiber.MethodHead {
		c.Request().ResetBody()
		// Fasthttp should skipbody by default if HEAD?
		c.Response().SkipBody = true
		c.Response().Header.SetContentLength(contentLength)
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}

	return nil
}
