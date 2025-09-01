package docs

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// APIDoc represents the API documentation structure
type APIDoc struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Servers    []Server            `json:"servers"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
	Tags       []Tag               `json:"tags"`
}

// Info represents API information
type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// Server represents server information
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// PathItem represents a path in the API
type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
}

// Operation represents an API operation
type Operation struct {
	Tags        []string              `json:"tags"`
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	OperationID string                `json:"operationId"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
}

// Parameter represents an API parameter
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Schema      *Schema     `json:"schema"`
	Example     interface{} `json:"example,omitempty"`
}

// RequestBody represents a request body
type RequestBody struct {
	Description string               `json:"description"`
	Required    bool                 `json:"required"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType represents media type content
type MediaType struct {
	Schema   *Schema            `json:"schema"`
	Example  interface{}        `json:"example,omitempty"`
	Examples map[string]Example `json:"examples,omitempty"`
}

// Response represents an API response
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Headers     map[string]Header    `json:"headers,omitempty"`
}

// Header represents a response header
type Header struct {
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}

// Schema represents a data schema
type Schema struct {
	Type        string            `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Description string            `json:"description,omitempty"`
	Properties  map[string]Schema `json:"properties,omitempty"`
	Required    []string          `json:"required,omitempty"`
	Items       *Schema           `json:"items,omitempty"`
	Example     interface{}       `json:"example,omitempty"`
	Ref         string            `json:"$ref,omitempty"`
}

// Example represents an example value
type Example struct {
	Summary       string      `json:"summary"`
	Description   string      `json:"description"`
	Value         interface{} `json:"value"`
	ExternalValue string      `json:"externalValue,omitempty"`
}

// Components represents reusable components
type Components struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

// SecurityScheme represents a security scheme
type SecurityScheme struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Name        string `json:"name,omitempty"`
	In          string `json:"in,omitempty"`
	Scheme      string `json:"scheme,omitempty"`
}

// Tag represents an API tag
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DocGenerator generates API documentation
type DocGenerator struct {
	docs    *APIDoc
	baseDir string
}

// NewDocGenerator creates a new documentation generator
func NewDocGenerator(baseDir string) *DocGenerator {
	return &DocGenerator{
		docs: &APIDoc{
			OpenAPI: "3.0.0",
			Info: Info{
				Title:       "Study Platform API",
				Description: "API documentation for the Study Platform backend services",
				Version:     "1.0.0",
			},
			Servers: []Server{
				{
					URL:         "http://localhost:8080",
					Description: "Development server",
				},
			},
			Paths:      make(map[string]PathItem),
			Components: Components{},
			Tags: []Tag{
				{Name: "auth", Description: "Authentication operations"},
				{Name: "users", Description: "User management operations"},
				{Name: "rooms", Description: "Room management operations"},
				{Name: "sessions", Description: "Study session operations"},
				{Name: "materials", Description: "Study material operations"},
				{Name: "friends", Description: "Friend management operations"},
				{Name: "realtime", Description: "Real-time communication operations"},
			},
		},
		baseDir: baseDir,
	}
}

// GenerateDocs generates documentation from Go source files
func (dg *DocGenerator) GenerateDocs() error {
	// Parse main.go to find route registrations
	if err := dg.parseMainFile(); err != nil {
		return fmt.Errorf("failed to parse main.go: %w", err)
	}

	// Parse handler files to extract endpoint information
	if err := dg.parseHandlerFiles(); err != nil {
		return fmt.Errorf("failed to parse handler files: %w", err)
	}

	// Generate OpenAPI JSON
	if err := dg.generateOpenAPIJSON(); err != nil {
		return fmt.Errorf("failed to generate OpenAPI JSON: %w", err)
	}

	// Generate HTML documentation
	if err := dg.generateHTMLDocs(); err != nil {
		return fmt.Errorf("failed to generate HTML docs: %w", err)
	}

	return nil
}

// parseMainFile parses main.go to find route registrations
func (dg *DocGenerator) parseMainFile() error {
	mainFile := filepath.Join(dg.baseDir, "cmd", "api", "main.go")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		return fmt.Errorf("main.go not found at %s", mainFile)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, mainFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse main.go: %w", err)
	}

	// Look for route registrations
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if fun.Sel.Name == "GET" || fun.Sel.Name == "POST" ||
					fun.Sel.Name == "PUT" || fun.Sel.Name == "DELETE" {
					dg.extractRouteInfo(call, fun.Sel.Name)
				}
			}
		}
		return true
	})

	return nil
}

// extractRouteInfo extracts route information from AST nodes
func (dg *DocGenerator) extractRouteInfo(call *ast.CallExpr, method string) {
	if len(call.Args) < 2 {
		return
	}

	// Extract path from first argument
	path := dg.extractStringLiteral(call.Args[0])
	if path == "" {
		return
	}

	// Extract handler name from second argument
	handlerName := dg.extractHandlerName(call.Args[1])
	if handlerName == "" {
		return
	}

	// Create or update path item
	pathItem, exists := dg.docs.Paths[path]
	if !exists {
		pathItem = PathItem{}
	}

	// Create operation
	operation := &Operation{
		Tags:        dg.extractTags(path),
		Summary:     dg.generateSummary(method, path),
		Description: dg.generateDescription(method, path),
		OperationID: fmt.Sprintf("%s%s", strings.ToLower(method), strings.ReplaceAll(path, "/", "")),
		Responses:   dg.generateDefaultResponses(),
	}

	// Add method to path item
	switch method {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "DELETE":
		pathItem.Delete = operation
	}

	dg.docs.Paths[path] = pathItem
}

// extractStringLiteral extracts string literal from AST node
func (dg *DocGenerator) extractStringLiteral(node ast.Expr) string {
	if lit, ok := node.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, `"`)
	}
	return ""
}

// extractHandlerName extracts handler function name from AST node
func (dg *DocGenerator) extractHandlerName(node ast.Expr) string {
	if ident, ok := node.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// extractTags extracts tags based on path
func (dg *DocGenerator) extractTags(path string) []string {
	if strings.Contains(path, "/auth") {
		return []string{"auth"}
	} else if strings.Contains(path, "/users") {
		return []string{"users"}
	} else if strings.Contains(path, "/rooms") {
		return []string{"rooms"}
	} else if strings.Contains(path, "/sessions") {
		return []string{"sessions"}
	} else if strings.Contains(path, "/materials") {
		return []string{"materials"}
	} else if strings.Contains(path, "/friends") {
		return []string{"friends"}
	} else if strings.Contains(path, "/realtime") {
		return []string{"realtime"}
	}
	return []string{"general"}
}

// generateSummary generates operation summary
func (dg *DocGenerator) generateSummary(method, path string) string {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < 2 {
		return fmt.Sprintf("%s %s", method, path)
	}

	resource := pathParts[1]
	switch method {
	case "GET":
		if strings.HasSuffix(path, "/") {
			return fmt.Sprintf("List %s", resource)
		}
		return fmt.Sprintf("Get %s", resource)
	case "POST":
		return fmt.Sprintf("Create %s", resource)
	case "PUT":
		return fmt.Sprintf("Update %s", resource)
	case "DELETE":
		return fmt.Sprintf("Delete %s", resource)
	default:
		return fmt.Sprintf("%s %s", method, resource)
	}
}

// generateDescription generates operation description
func (dg *DocGenerator) generateDescription(method, path string) string {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < 2 {
		return fmt.Sprintf("Perform %s operation on %s", method, path)
	}

	resource := pathParts[1]
	switch method {
	case "GET":
		if strings.HasSuffix(path, "/") {
			return fmt.Sprintf("Retrieve a list of %s", resource)
		}
		return fmt.Sprintf("Retrieve a specific %s by ID", resource)
	case "POST":
		return fmt.Sprintf("Create a new %s", resource)
	case "PUT":
		return fmt.Sprintf("Update an existing %s", resource)
	case "DELETE":
		return fmt.Sprintf("Delete a %s", resource)
	default:
		return fmt.Sprintf("Perform %s operation on %s", method, resource)
	}
}

// generateDefaultResponses generates default API responses
func (dg *DocGenerator) generateDefaultResponses() map[string]Response {
	return map[string]Response{
		"200": {
			Description: "Successful operation",
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Type: "object",
						Properties: map[string]Schema{
							"status": {Type: "string"},
							"data":   {Type: "object"},
						},
					},
				},
			},
		},
		"400": {
			Description: "Bad request",
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Type: "object",
						Properties: map[string]Schema{
							"error":   {Type: "string"},
							"message": {Type: "string"},
						},
					},
				},
			},
		},
		"401": {
			Description: "Unauthorized",
		},
		"500": {
			Description: "Internal server error",
		},
	}
}

// parseHandlerFiles parses handler files to extract endpoint information
func (dg *DocGenerator) parseHandlerFiles() error {
	handlerDir := filepath.Join(dg.baseDir, "internal")
	return filepath.Walk(handlerDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_handler.go") {
			if err := dg.parseHandlerFile(path); err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}
		}
		return nil
	})
}

// parseHandlerFile parses a single handler file
func (dg *DocGenerator) parseHandlerFile(filePath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Extract function documentation from comments
	ast.Inspect(node, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			dg.extractFunctionDocs(funcDecl)
		}
		return true
	})

	return nil
}

// extractFunctionDocs extracts documentation from function declarations
func (dg *DocGenerator) extractFunctionDocs(funcDecl *ast.FuncDecl) {
	if funcDecl.Doc == nil {
		return
	}

	// Extract comments and look for API documentation
	for _, comment := range funcDecl.Doc.List {
		text := comment.Text
		if strings.HasPrefix(text, "// @") {
			dg.parseAPIDocComment(text, funcDecl.Name.Name)
		}
	}
}

// parseAPIDocComment parses API documentation comments
func (dg *DocGenerator) parseAPIDocComment(comment, funcName string) {
	// Simple parsing of @api comments
	// In a full implementation, you'd want more sophisticated parsing
	if strings.Contains(comment, "@api") {
		// Extract API information from comment
		// This is a simplified version
	}
}

// generateOpenAPIJSON generates OpenAPI JSON file
func (dg *DocGenerator) generateOpenAPIJSON() error {
	outputFile := filepath.Join(dg.baseDir, "docs", "openapi.json")

	// Ensure docs directory exists
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(dg.docs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal OpenAPI spec: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write OpenAPI spec: %w", err)
	}

	fmt.Printf("Generated OpenAPI spec: %s\n", outputFile)
	return nil
}

// generateHTMLDocs generates HTML documentation
func (dg *DocGenerator) generateHTMLDocs() error {
	outputFile := filepath.Join(dg.baseDir, "docs", "index.html")

	// HTML template for API documentation
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - API Documentation</title>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
        .swagger-ui { font-family: Arial, sans-serif; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: './openapi.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	// Parse template
	tmpl, err := template.New("api-docs").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer file.Close()

	// Execute template
	data := struct {
		Title string
	}{
		Title: dg.docs.Info.Title,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	fmt.Printf("Generated HTML docs: %s\n", outputFile)
	return nil
}
