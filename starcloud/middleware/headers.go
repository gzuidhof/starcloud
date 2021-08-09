/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package middleware

import "github.com/gofiber/fiber/v2"

func AddCommonCDNHeadersMiddleware(ctx *fiber.Ctx) error {
	err := ctx.Next()

	ctx.Set(fiber.HeaderAccessControlAllowOrigin, "*")
	if (ctx.Secure()) {
		ctx.Set(fiber.HeaderStrictTransportSecurity, "max-age=31536000; includeSubDomains; preload")
	}
	ctx.Set(fiber.HeaderXContentTypeOptions, "nosniff")
	ctx.Set(fiber.HeaderXXSSProtection, "1; mode=block")
	
	ctx.Set("Cross-Origin-Embedder-Policy", "require-corp")
	ctx.Set("Cross-Origin-Opener-Policy", "same-origin")

	return err
}
