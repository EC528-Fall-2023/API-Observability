package parserimport

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var frameworkSignatures = map[string]string{
	"github.com/gin-gonic/gin":    "Gin",
	"github.com/labstack/echo/v4": "Echo",
	"github.com/gorilla/mux":      "GorillaMux",
	"net/http":                    "net/http",
	"github.com/gofiber/fiber/v2": "Fiber",
}

// Endpoint represents a single API endpoint
type Endpoint struct {
	Method  string
	Pattern string
	Handler string
}

// DetectFrameworkAndEndpoints will analyze the provided source file and return the detected framework and endpoints
func DetectFrameworkAndEndpoints(sourceFilePath string) (string, []Endpoint, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, sourceFilePath, nil, 0)
	if err != nil {
		return "", nil, err
	}

	detectedFramework, isNetHttpDetected := detectFramework(f)

	if detectedFramework == "" && isNetHttpDetected {
		detectedFramework = "net/http"
	} else if detectedFramework == "" {
		return "", nil, fmt.Errorf("no known framework detected")
	}

	endpoints := parseForEndpoints(f, detectedFramework)
	return detectedFramework, endpoints, nil
}

func detectFramework(f *ast.File) (string, bool) {
	var detectedFramework string
	var isNetHttpDetected bool

	ast.Inspect(f, func(n ast.Node) bool {
		importSpec, ok := n.(*ast.ImportSpec)
		if !ok {
			return true
		}

		if importSpec.Path != nil {
			for signature, framework := range frameworkSignatures {
				if importSpec.Path.Value == `"`+signature+`"` {
					if framework == "net/http" {
						isNetHttpDetected = true
					} else {
						detectedFramework = framework
						return false // Stop if a specific framework is found
					}
				}
			}
		}
		return true
	})

	return detectedFramework, isNetHttpDetected
}

func parseForEndpoints(f *ast.File, framework string) []Endpoint {
	var endpoints []Endpoint
	switch framework {
	case "Gin":
		endpoints = parseGinEndpoints(f)
	case "Echo":
		endpoints = parseEchoEndpoints(f)
	case "GorillaMux":
		endpoints = parseGorillaMuxEndpoints(f)
	case "net/http":
		endpoints = parseNetHttpEndpoints(f)
	case "Fiber":
		endpoints = parseFiberEndpoints(f)
	}

	return endpoints
}

// Add this function to your existing code

func parseGinEndpoints(f *ast.File) []Endpoint {
	var endpoints []Endpoint
	ast.Inspect(f, func(n ast.Node) bool {
		// Look for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check for HTTP method calls (GET, POST, etc.)
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && isGinHTTPMethod(selExpr.Sel.Name) {
			// Assuming the first argument is the URL pattern
			if len(callExpr.Args) >= 2 {
				pattern := ""
				handler := ""

				// Get pattern
				if patternLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && patternLit.Kind == token.STRING {
					pattern = strings.Trim(patternLit.Value, `"`)
				}

				// Get handler
				if funcLit, ok := callExpr.Args[1].(*ast.Ident); ok {
					handler = funcLit.Name
				} else if funcLit, ok := callExpr.Args[1].(*ast.SelectorExpr); ok {
					// This handles cases where the handler is a method of an object
					handler = fmt.Sprintf("%s.%s", funcLit.X, funcLit.Sel.Name)
				}

				if pattern != "" && handler != "" {
					endpoints = append(endpoints, Endpoint{
						Method:  selExpr.Sel.Name,
						Pattern: pattern,
						Handler: handler, // Add the handler to your Endpoint struct
					})
				}
			}
		}
		return true
	})

	return endpoints
}

func isGinHTTPMethod(methodName string) bool {
	switch methodName {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS":
		return true
	}
	return false
}

// parseEchoEndpoints extracts the API endpoints from an Echo based Go source file.
func parseEchoEndpoints(f *ast.File) []Endpoint {
	var endpoints []Endpoint

	ast.Inspect(f, func(n ast.Node) bool {
		// Look for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a method call (e.g., e.GET)
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && isEchoHTTPMethod(selExpr.Sel.Name) {
			// Assuming the first argument is the URL pattern
			if len(callExpr.Args) >= 2 {
				pattern := ""
				handler := ""

				// Get pattern
				if patternLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && patternLit.Kind == token.STRING {
					pattern = strings.Trim(patternLit.Value, `"`)
				}

				// Get handler
				if funcLit, ok := callExpr.Args[1].(*ast.Ident); ok {
					handler = funcLit.Name
				} else if funcLit, ok := callExpr.Args[1].(*ast.SelectorExpr); ok {
					// This handles cases where the handler is a method of an object
					handler = fmt.Sprintf("%s.%s", funcLit.X, funcLit.Sel.Name)
				}

				if pattern != "" && handler != "" {
					endpoints = append(endpoints, Endpoint{
						Method:  selExpr.Sel.Name,
						Pattern: pattern,
						Handler: handler, // Add the handler to your Endpoint struct
					})
				}
			}
		}
		return true
	})

	return endpoints
}

// isEchoHTTPMethod checks if a method name is one of the HTTP methods used in Echo
func isEchoHTTPMethod(methodName string) bool {
	switch methodName {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS":
		return true
	}
	return false
}

// parseGorillaMuxEndpoints extracts the API endpoints from a Gorilla Mux based Go source file.
func parseGorillaMuxEndpoints(f *ast.File) []Endpoint {
	var endpoints []Endpoint

	ast.Inspect(f, func(n ast.Node) bool {
		// Look for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a HandleFunc call
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && selExpr.Sel.Name == "HandleFunc" {
			// Assuming the first argument is the URL pattern
			if len(callExpr.Args) >= 2 {
				pattern := ""
				handler := ""

				// Get pattern
				if patternLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && patternLit.Kind == token.STRING {
					pattern = strings.Trim(patternLit.Value, `"`)
				}

				// Get handler
				if funcLit, ok := callExpr.Args[1].(*ast.Ident); ok {
					handler = funcLit.Name
				} else if funcLit, ok := callExpr.Args[1].(*ast.SelectorExpr); ok {
					// This handles cases where the handler is a method of an object
					handler = fmt.Sprintf("%s.%s", funcLit.X, funcLit.Sel.Name)
				}

				if pattern != "" && handler != "" {
					endpoints = append(endpoints, Endpoint{
						Method:  "CUSTOM", // Gorilla Mux usually uses custom methods for HandleFunc
						Pattern: pattern,
						Handler: handler,
					})
				}
			}
		}

		return true
	})

	return endpoints
}

// parseNetHttpEndpoints extracts the API endpoints from a Go source file using the net/http package.
func parseNetHttpEndpoints(f *ast.File) []Endpoint {
	var endpoints []Endpoint

	ast.Inspect(f, func(n ast.Node) bool {
		// Look for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a HandleFunc or Handle call
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && (selExpr.Sel.Name == "HandleFunc" || selExpr.Sel.Name == "Handle") {
			// Assuming the first argument is the URL pattern
			if len(callExpr.Args) >= 2 {
				pattern := ""
				handler := ""

				// Get pattern
				if patternLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && patternLit.Kind == token.STRING {
					pattern = strings.Trim(patternLit.Value, `"`)
				}

				// Get handler
				if funcLit, ok := callExpr.Args[1].(*ast.Ident); ok {
					handler = funcLit.Name
				} else if funcLit, ok := callExpr.Args[1].(*ast.SelectorExpr); ok {
					// This handles cases where the handler is a method of an object
					handler = fmt.Sprintf("%s.%s", funcLit.X, funcLit.Sel.Name)
				}

				if pattern != "" && handler != "" {
					method := "ALL"
					if selExpr.Sel.Name == "Handle" {
						method = "CUSTOM"
					}
					endpoints = append(endpoints, Endpoint{
						Method:  method,
						Pattern: pattern,
						Handler: handler,
					})
				}
			}
		}

		return true
	})

	return endpoints
}

// parseFiberEndpoints extracts the API endpoints from a Fiber based Go source file.
func parseFiberEndpoints(f *ast.File) []Endpoint {
	var endpoints []Endpoint

	ast.Inspect(f, func(n ast.Node) bool {
		// Look for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a Fiber method call (e.g., app.Get)
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && isFiberHTTPMethod(selExpr.Sel.Name) {
			// Assuming the first argument is the URL pattern
			if len(callExpr.Args) >= 2 {
				pattern := ""
				handler := ""

				// Get pattern
				if patternLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && patternLit.Kind == token.STRING {
					pattern = strings.Trim(patternLit.Value, `"`)
				}

				// Get handler
				if funcLit, ok := callExpr.Args[1].(*ast.Ident); ok {
					handler = funcLit.Name
				} else if funcLit, ok := callExpr.Args[1].(*ast.SelectorExpr); ok {
					// This handles cases where the handler is a method of an object
					handler = fmt.Sprintf("%s.%s", funcLit.X, funcLit.Sel.Name)
				}

				if pattern != "" && handler != "" {
					endpoints = append(endpoints, Endpoint{
						Method:  strings.ToUpper(selExpr.Sel.Name), // Convert method to uppercase
						Pattern: pattern,
						Handler: handler,
					})
				}
			}
		}
		return true
	})

	return endpoints
}

// isFiberHTTPMethod checks if a method name is one of the HTTP methods used in Fiber
func isFiberHTTPMethod(methodName string) bool {
	// Fiber uses lowercase method names (Get, Post, etc.)
	switch strings.ToLower(methodName) {
	case "get", "post", "put", "delete", "patch", "options":
		return true
	}
	return false
}
