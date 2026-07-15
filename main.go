package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

//go:embed template.exe
var templateBytes []byte

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
					speechLabel.SetText("Stoking the fire... Preparing a delicious meal!")
					statusLabel.SetText("Baking a valid Windows executable...")
					progressBar.SetValue(80)
				})
				time.Sleep(800 * time.Millisecond)

				meal := cookMealFromPaths(files)

				// Cook the meal by generating the executable
				outputPath, err := generateMealExe(meal)

				mainWindow.Synchronize(func() {
					progressBar.SetValue(100)
					if err != nil {
						speechLabel.SetText("Oops! I burnt the files: " + err.Error())
						statusLabel.SetText("Cooking failed!")
					} else {
						speechLabel.SetText(fmt.Sprintf("Voila! Behold your %s! It has been compiled into a real EXE on your Desktop/working directory! *ur mama eating sfx*", meal.Name))
						statusLabel.SetText(fmt.Sprintf("Ready! Saved as %s", filepath.Base(outputPath)))
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

func generateMealExe(meal MealInfo) (string, error) {
	parts := strings.Split(meal.Name, " ")
	baseType := parts[len(parts)-1]
	asciiArt, exists := asciiArtMap[baseType]
	if !exists {
		asciiArt = asciiArtMap["Default"]
	}

	// Payload formatting
	magicSignature := "B0BBY_CHEF_REC1PE_START"
	recipeList := strings.Join(meal.Recipe, "\n")
	payload := fmt.Sprintf("%s%s|||%s|||%s|||%s", magicSignature, meal.Name, meal.Description, recipeList, asciiArt)

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

	// Read precompiled template from embed, append payload, and save!
	var finalBytes []byte
	if len(templateBytes) > 0 {
		finalBytes = make([]byte, len(templateBytes)+len(payload))
		copy(finalBytes, templateBytes)
		copy(finalBytes[len(templateBytes):], []byte(payload))
	} else {
		return "", fmt.Errorf("embedded template.exe is empty")
	}

	if err := os.WriteFile(outputPath, finalBytes, 0755); err != nil {
		return "", err
	}

	return outputPath, nil
}
