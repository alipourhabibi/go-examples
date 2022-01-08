package api

// swagger:parameters loginUserParameter
type loginUserWrapper struct {
	// in: body
	Body LogInUser
}

// swagger:parameters registerUserParameter
type registerUserWrapper struct {
	// in: body
	Body RegisterUser
}

// swagger:parameters logOutUserParameter
type logoutUserWrapper struct {
	// in: header
	// name: Autorization
	// Example: Bearer TOKEN
	// Required: true
	Authorization string `json:"authorization"`
}

// swagger:parameters refreshTokenParameter
type refreshTokenWrapper struct {
	// in: body
	Body struct {
		// Example: TOKEN
		// Required: true
		RefreshToken string `json:"refresh_token"`
	}
}

// swagger:parameters newPostParameter
type newPostWrapper struct {
	// in: header
	// name: Autorization
	// Example: Bearer TOKEN
	// Required: true
	Authorization string `json:"authorization"`
	// in: body
	Body Post
}

// swagger:parameters deletePostParameter
type deletePostWrapper struct {
	// in: header
	// name: Autorization
	// Example: Bearer TOKEN
	// Required: true
	Authorization string `json:"authorization"`
	// in: path
	ID string `json:"id"`
	// in: body
	Body Post
}

// swagger:parameters updatePostParameter
type updatePostWrapper struct {
	// in: header
	// name: Autorization
	// Example: Bearer TOKEN
	// Required: true
	Authorization string `json:"authorization"`
	// in: path
	ID string `json:"id"`
	// in: body
	Body Post
}

// swagger:parameters getPostParameter
type getPostWrapper struct {
	// in: path
	ID string `json:"id"`
}

// swagger:response loginSuccess
type loginSuccessWrapper struct {
	// in: body
	Body struct {
		// Example: TOKEN
		AccessToken string `json:"access_token"`
		// Example: TOKEN
		RefreshToken string `json:"refresh_token"`
	}
}

// swagger:response responseBadRequest
type responseBadRequestWrapper struct {
	// in: body
	Body struct {
		// Example: Invalid JSON provided
		MSG string `json:"msg"`
	}
}

// swagger:response responseUnauthorized
type responseUnauthorizedWrapper struct {
	// in: body
	Body struct {
		// Example: Unauthorized
		MSG string `json:"msg"`
	}
}

// swagger:response responseInternalServerError
type responseInternalServerErrorWrapper struct {
	// in: body
	Body struct {
		// Example: Internal server error
		MSG string `json:"msg"`
	}
}

// swagger:response responseCreated
type responseCreatedWrapper struct {
	// in: body
	Body struct {
		// Example: created
		MSG string `json:"msg"`
	}
}

// swagger:response responseSuccess
type responseSuccessWrapper struct {
	// in: body
	Body struct {
		// Example: Success
		MSG string `json:"msg"`
	}
}

// swagger:response responseGetDataSuccess
type responseGetDataSuccess struct {
	// in: body
	Dody struct {
		Datas struct {
			// Example: 1
			ID int `json:"id"`
			// Example: ali
			Username string `json:"username"`
			// Example: title
			Title string `json:"title"`
			// Example: text
			Text string `json:"text"`
		} `json:"datas"`
		// Example: Success
		MSG string `json:"msg"`
	}
}
