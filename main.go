package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MealType struct {
	Keywords []string
	Type     string
	Suffix   string
}

var mealTypes = []MealType{
	{Keywords: []string{"png", "jpg", "jpeg", "gif", "bmp"}, Type: "Stew", Suffix: "with visual garnish"},
	{Keywords: []string{"txt", "md", "pdf", "doc", "docx"}, Type: "Soup", Suffix: "infused with literary wisdom"},
	{Keywords: []string{"mp3", "wav", "ogg", "flac", "m4a"}, Type: "Smoothie", Suffix: "with melodic undertones"},
	{Keywords: []string{"mp4", "mkv", "avi", "mov", "webm"}, Type: "Casserole", Suffix: "garnished with visual action"},
	{Keywords: []string{"js", "py", "go", "rs", "cpp", "c", "java", "html", "css", "json"}, Type: "Stir-Fry", Suffix: "with intellectual flavor"},
	{Keywords: []string{"exe", "msi", "bin", "sh", "bat"}, Type: "Roast", Suffix: "infused with high-energy instructions"},
	{Keywords: []string{"zip", "tar", "gz", "rar", "7z"}, Type: "Pie", Suffix: "tightly packed with rich goodness"},
}

var fallbackMeals = []string{"Sauté", "Curry", "Risotto", "Gumbo", "Paella", "Chowder", "Fricassee"}
var fallbackGarnishes = []string{"with a dash of mystery", "cooked to absolute perfection", "served extra hot", "that tastes surprisingly good"}

var asciiArtMap = map[string]string{
	"Soup": `
      (  )   (   )  )
       ) (   )  (  (
       ___________
      (___________)
       |         |
       |  SOUP   |
       \_________/
  `,
	"Stew": `
         (  (
          )  )
       ___________
      (___________)
       \         /
        \_STEW__/
  `,
	"Smoothie": `
         \_|_/
          \ /
         [===]
         |   |
         | * |
         |___|
         (===)
  `,
	"Casserole": `
     =================
    | [][][][][][][]  |
    |  CASSEROLE      |
    |                 |
     =================
  `,
	"Stir-Fry": `
       (   )  (   )
        \_/_/_/_/_/
        |         |
        |  FRY!   |
        \_________/
  `,
	"Roast": `
       ,---.  ,---.
      /     \/     \
      \            /
       '--.    .--'
          |    |
          |____|
  `,
	"Pie": `
         ______
       .-'      '-.
     .'  _|_|_|_   '.
     |  |_|_|_|_|   |
     '.   PIE     .'
       '-.______.-'
  `,
	"Default": `
       (\___/)
       (='.'=)
       (")_(")
       Chef Bobby's
       Secret Recipe!
  `,
}

type MealInfo struct {
	Name        string
	Description string
	Recipe      []string
}

func main() {
	var mainWindow *walk.MainWindow
	var statusLabel *walk.Label
	var speechLabel *walk.Label
	var progressBar *walk.ProgressBar

	if err := (MainWindow{
		AssignTo: &mainWindow,
		Title:    "Chef Bobby's Cooking Game",
		MinSize:  Size{600, 450},
		Layout:   VBox{},
		OnDropFiles: func(files []string) {
			if len(files) == 0 {
				return
			}

			progressBar.SetVisible(true)
			progressBar.SetValue(0)

			go func() {
				// We use a small utility to safety set text in UI thread
				mainWindow.Synchronize(func() {
					speechLabel.SetText("Oh, yes! I'm chopping up the raw file data...")
					statusLabel.SetText("Dicing binary blocks...")
					progressBar.SetValue(20)
				})
				time.Sleep(800 * time.Millisecond)

				mainWindow.Synchronize(func() {
					speechLabel.SetText("Stirring the pot! Let the headers simmer nicely...")
					statusLabel.SetText("Boiling file headers...")
					progressBar.SetValue(50)
				})
				time.Sleep(800 * time.Millisecond)

				mainWindow.Synchronize(func() {
					speechLabel.SetText("Stoking the compiler fire... Creating a fine Windows meal!")
					statusLabel.SetText("Baking a valid x86-64 Windows executable...")
					progressBar.SetValue(80)
				})
				time.Sleep(800 * time.Millisecond)

				// Analyze dropped files to cook recipe
				meal := cookMealFromPaths(files)

				// Compile x86-64 executable
				err := compileMealExe(meal)

				mainWindow.Synchronize(func() {
					progressBar.SetValue(100)
					if err != nil {
						speechLabel.SetText("Oops! I burnt the files: " + err.Error())
						statusLabel.SetText("Cooking failed!")
					} else {
						speechLabel.SetText(fmt.Sprintf("Voila! Behold your %s! It has been compiled into a real EXE on your Desktop/working directory! *ur mama eating sfx*", meal.Name))
						statusLabel.SetText("Ready to serve!")
					}
				})
			}()
		},
		Children: []Widget{
			Label{
				Text:          "Chef Bobby's Kitchen",
				Font:          Font{PointSize: 16, Bold: true},
				TextAlignment: AlignCenter,
			},
			Label{
				Text:          "Drag files into this window, and Chef Bobby will cook them into a real, valid Windows EXE meal!",
				TextAlignment: AlignCenter,
			},
			Label{
				AssignTo:      &speechLabel,
				Text:          "Welcome to my kitchen! Drag your files here, and let's get cooking!",
				TextAlignment: AlignCenter,
				Font:          Font{PointSize: 11, Italic: true},
			},
			ProgressBar{
				AssignTo: &progressBar,
				MinValue: 0,
				MaxValue: 100,
				Visible:  false,
			},
			Label{
				AssignTo:      &statusLabel,
				Text:          "Stove is off. Ready for file ingredients.",
				TextAlignment: AlignCenter,
			},
		},
	}.Create()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mainWindow.Run()
}

func cookMealFromPaths(paths []string) MealInfo {
	var totalSize int64
	var extensionCounts = make(map[string]int)

	for _, path := range paths {
		ext := strings.ToLower(filepath.Ext(path))
		ext = strings.TrimPrefix(ext, ".")
		if ext != "" {
			extensionCounts[ext]++
		}
		if fi, err := os.Stat(path); err == nil {
			totalSize += fi.Size()
		}
	}

	maxExt := ""
	maxCount := 0
	for ext, count := range extensionCounts {
		if count > maxCount {
			maxCount = count
			maxExt = ext
		}
	}

	mealType := ""
	garnish := ""

	for _, m := range mealTypes {
		for _, kw := range m.Keywords {
			if kw == maxExt {
				mealType = m.Type
				garnish = m.Suffix
				break
			}
		}
		if mealType != "" {
			break
		}
	}

	if mealType == "" {
		h := int(totalSize) % len(fallbackMeals)
		if h < 0 {
			h = -h
		}
		mealType = fallbackMeals[h]
		garnish = fallbackGarnishes[int(totalSize)%len(fallbackGarnishes)]
	}

	mainIngredient := "Unknown File"
	if len(paths) > 0 {
		base := strings.TrimSuffix(filepath.Base(paths[0]), filepath.Ext(paths[0]))
		mainIngredient = base
		if len(paths) > 1 {
			mainIngredient += " & Friends"
		}
	}

	if len(mainIngredient) > 0 {
		mainIngredient = strings.Title(mainIngredient)
	}
	if len(mainIngredient) > 25 {
		mainIngredient = mainIngredient[:22] + "..."
	}

	mealName := fmt.Sprintf("%s %s", mainIngredient, mealType)
	sizeKB := float64(totalSize) / 1024.0

	description := fmt.Sprintf("A grand culinary masterpiece cooked by Chef Bobby! He combined %d file(s) (%.1f KB total) into a spectacular %s, %s. Warning: Ur mama might eat this too fast!", len(paths), sizeKB, mealName, garnish)

	var recipe []string
	for _, path := range paths {
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
		if ext == "" {
			ext = "data"
		}
		recipe = append(recipe, fmt.Sprintf("- 1 pinch of %s (pure %s essence)", filepath.Base(path), ext))
	}

	return MealInfo{
		Name:        mealName,
		Description: description,
		Recipe:      recipe,
	}
}

func compileMealExe(meal MealInfo) error {
	parts := strings.Split(meal.Name, " ")
	baseType := parts[len(parts)-1]
	asciiArt, exists := asciiArtMap[baseType]
	if !exists {
		asciiArt = asciiArtMap["Default"]
	}

	buildId := fmt.Sprintf("build_%d_%d", time.Now().Unix(), rand.Intn(1000))
	buildDir := filepath.Join(os.TempDir(), buildId)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(buildDir)

	var recipePrinters []string
	for _, r := range meal.Recipe {
		recipePrinters = append(recipePrinters, fmt.Sprintf("\tfmt.Println(%q)", r))
	}

	goCode := fmt.Sprintf(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("    YOU ARE PLAYING CHEF BOBBY'S COOKING GAME!    ")
	fmt.Println("==================================================")
	fmt.Println("")
	fmt.Println("Chef Bobby has prepared a magnificent executable meal for you:")
	fmt.Println("🍳 MEAL: %s")
	fmt.Println("")
	fmt.Println("--- INGREDIENTS USED ---")
%s
	fmt.Println("")
	fmt.Println("--- ANALYSIS ---")
	fmt.Println(%q)
	fmt.Println("")
	fmt.Println("--- MEAL PRESENTATION ---")
	fmt.Println(%q)
	fmt.Println("")
	fmt.Println("==================================================")
	fmt.Println("    *UR MAMA EATING SFX* (CRUNCH MUNCH GULP!)      ")
	fmt.Println("==================================================")
	fmt.Println("Press ENTER to finish digesting this amazing meal...")

	var input string
	fmt.Scanln(&input)
}
`, meal.Name, strings.Join(recipePrinters, "\n"), meal.Description, asciiArt)

	mainGoPath := filepath.Join(buildDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(goCode), 0644); err != nil {
		return err
	}

	cmdInit := exec.Command("go", "mod", "init", "bobby_meal")
	cmdInit.Dir = buildDir
	if err := cmdInit.Run(); err != nil {
		return err
	}

	cleanName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, strings.ToLower(meal.Name))

	exeFilename := cleanName + ".exe"
	workingDir, err := os.Getwd()
	if err != nil {
		workingDir = "."
	}

	homeDir, err := os.UserHomeDir()
	outputPath := filepath.Join(workingDir, exeFilename)
	if err == nil {
		desktopPath := filepath.Join(homeDir, "Desktop")
		if _, err := os.Stat(desktopPath); err == nil {
			outputPath = filepath.Join(desktopPath, exeFilename)
		}
	}

	cmdBuild := exec.Command("go", "build", "-ldflags", "-s -w", "-o", outputPath, "main.go")
	cmdBuild.Dir = buildDir
	cmdBuild.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64")

	if err := cmdBuild.Run(); err != nil {
		return err
	}

	return nil
}
