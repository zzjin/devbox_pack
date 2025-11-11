package providers

import (
	"strings"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewStaticFileProvider(t *testing.T) {
	provider := NewStaticFileProvider()

	if provider == nil {
		t.Fatal("NewStaticFileProvider() returned nil")
	}

	if provider.GetName() != "staticfile" {
		t.Errorf("expected name 'staticfile', got %s", provider.GetName())
	}

	if provider.GetPriority() != 90 {
		t.Errorf("expected priority 90, got %d", provider.GetPriority())
	}
}

func TestStaticFileProvider_Detect_NoStaticFiles(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "README.md", IsDirectory: false},
		{Path: "main.py", IsDirectory: false},
		{Path: "package.json", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if result.Matched {
		t.Error("expected not matched for non-static file project")
	}
}

func TestStaticFileProvider_Detect_WithHTMLFiles(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "index.html", IsDirectory: false},
		{Path: "about.html", IsDirectory: false},
		{Path: "styles.css", IsDirectory: false},
		{Path: "script.js", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with HTML files")
	}

	if result.Language != "staticfile" {
		t.Errorf("expected language 'staticfile', got %s", result.Language)
	}
}

func TestStaticFileProvider_Detect_WithIndexHTML(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "index.html", IsDirectory: false},
		{Path: "assets/", IsDirectory: true},
		{Path: "assets/style.css", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with index.html")
	}
}

func TestStaticFileProvider_Detect_WithStaticfile(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "Staticfile", IsDirectory: false},
		{Path: "public/", IsDirectory: true},
		{Path: "public/index.html", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with Staticfile")
	}
}

func TestStaticFileProvider_Detect_WithCSSAndJS(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "styles.css", IsDirectory: false},
		{Path: "main.css", IsDirectory: false},
		{Path: "app.js", IsDirectory: false},
		{Path: "utils.js", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with CSS and JS files")
	}
}

func TestStaticFileProvider_Detect_WithImages(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "index.html", IsDirectory: false},
		{Path: "logo.png", IsDirectory: false},
		{Path: "banner.jpg", IsDirectory: false},
		{Path: "icon.svg", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with images")
	}
}

func TestStaticFileProvider_Detect_WithAssetsDirectory(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "index.html", IsDirectory: false},
		{Path: "assets/", IsDirectory: true},
		{Path: "assets/css/", IsDirectory: true},
		{Path: "assets/js/", IsDirectory: true},
		{Path: "assets/images/", IsDirectory: true},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with assets directory")
	}
}

func TestStaticFileProvider_Detect_WithCommonFiles(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "index.html", IsDirectory: false},
		{Path: "favicon.ico", IsDirectory: false},
		{Path: "robots.txt", IsDirectory: false},
		{Path: "sitemap.xml", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with common web files")
	}
}

func TestStaticFileProvider_Detect_ExcludeOtherLanguages(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Test with Node.js project that has HTML files
	files := []types.FileInfo{
		{Path: "index.html", IsDirectory: false},
		{Path: "package.json", IsDirectory: false},
		{Path: "node_modules/", IsDirectory: true},
		{Path: "src/", IsDirectory: true},
		{Path: "src/app.js", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	// Should have lower confidence due to presence of package.json
	if result.Matched && result.Confidence > 0.7 {
		t.Error("expected lower confidence for project with other language indicators")
	}
}

func TestStaticFileProvider_Detect_ExcludePythonProject(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Test with Python project that has HTML templates
	files := []types.FileInfo{
		{Path: "templates/", IsDirectory: true},
		{Path: "templates/index.html", IsDirectory: false},
		{Path: "requirements.txt", IsDirectory: false},
		{Path: "app.py", IsDirectory: false},
		{Path: "static/", IsDirectory: true},
		{Path: "static/style.css", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	// Should have lower confidence due to presence of Python files
	if result.Matched && result.Confidence > 0.7 {
		t.Error("expected lower confidence for Python project with HTML templates")
	}
}

func TestStaticFileProvider_GenerateCommands_BasicProject(t *testing.T) {
	provider := NewStaticFileProvider()

	result := &types.DetectResult{
		Matched:  true,
		Language: "staticfile",
		Evidence: types.Evidence{
			Files: []string{"index.html", "styles.css", "script.js"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Run) == 0 {
		t.Error("expected run commands")
	}

	// Check for static file server command (nginx is the actual implementation)
	foundServerCommand := false
	for _, cmd := range commands.Run {
		if strings.Contains(cmd, "nginx") {
			foundServerCommand = true
			break
		}
	}
	if !foundServerCommand {
		t.Error("expected nginx server command in run commands")
	}
}

func TestStaticFileProvider_GenerateCommands_WithStaticfile(t *testing.T) {
	provider := NewStaticFileProvider()

	result := &types.DetectResult{
		Matched:  true,
		Language: "staticfile",
		Evidence: types.Evidence{
			Files: []string{"Staticfile", "public/index.html"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Run) == 0 {
		t.Error("expected run commands")
	}

	// Should handle Staticfile configuration
	foundCommand := false
	for _, cmd := range commands.Run {
		if len(cmd) > 0 {
			foundCommand = true
			break
		}
	}
	if !foundCommand {
		t.Error("expected at least one run command")
	}
}

func TestStaticFileProvider_GenerateEnvironment(t *testing.T) {
	provider := NewStaticFileProvider()

	result := &types.DetectResult{
		Matched:  true,
		Language: "staticfile",
	}

	env := provider.GenerateEnvironment(result)

	if env["NGINX_PORT"] != "80" {
		t.Errorf("expected NGINX_PORT '80', got %s", env["NGINX_PORT"])
	}

	if env["DOCUMENT_ROOT"] != "/usr/share/nginx/html" {
		t.Errorf("expected DOCUMENT_ROOT '/usr/share/nginx/html', got %s", env["DOCUMENT_ROOT"])
	}
}

func TestStaticFileProvider_NeedsNativeCompilation(t *testing.T) {
	provider := NewStaticFileProvider()

	result := &types.DetectResult{
		Matched:  true,
		Language: "staticfile",
	}

	needsCompilation := provider.NeedsNativeCompilation(result)
	if needsCompilation {
		t.Error("expected static file projects to not need native compilation")
	}
}

func TestStaticFileProvider_Detect_HTMExtension(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "index.htm", IsDirectory: false},
		{Path: "about.htm", IsDirectory: false},
		{Path: "contact.html", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with .htm files")
	}
}

func TestStaticFileProvider_Detect_PublicDirectory(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "public/", IsDirectory: true},
		{Path: "public/index.html", IsDirectory: false},
		{Path: "public/css/", IsDirectory: true},
		{Path: "public/js/", IsDirectory: true},
		{Path: "static/", IsDirectory: true},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with public directory")
	}
}

func TestStaticFileProvider_Detect_ImageFormats(t *testing.T) {
	provider := NewStaticFileProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "index.html", IsDirectory: false},
		{Path: "image.png", IsDirectory: false},
		{Path: "photo.jpeg", IsDirectory: false},
		{Path: "animation.gif", IsDirectory: false},
		{Path: "vector.svg", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for static file project with various image formats")
	}
}
