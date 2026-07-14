const express = require('express');
const multer = require('multer');
const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3000;

// Setup uploads
const upload = multer({ dest: '/tmp/chef-uploads/' });

app.use(express.static('public'));
app.use(express.json());

// List of possible meal titles and suffixes
const mealTypes = [
  { keywords: ['img', 'png', 'jpg', 'jpeg', 'gif', 'bmp'], type: 'Stew', suffix: 'with visual garnish' },
  { keywords: ['txt', 'md', 'pdf', 'doc', 'docx'], type: 'Soup', suffix: 'infused with literary wisdom' },
  { keywords: ['mp3', 'wav', 'ogg', 'flac', 'm4a'], type: 'Smoothie', suffix: 'with melodic undertones' },
  { keywords: ['mp4', 'mkv', 'avi', 'mov', 'webm'], type: 'Casserole', suffix: 'garnished with visual action' },
  { keywords: ['js', 'py', 'go', 'rs', 'cpp', 'c', 'java', 'html', 'css', 'json'], type: 'Stir-Fry', suffix: 'with intellectual flavor' },
  { keywords: ['exe', 'msi', 'bin', 'sh', 'bat'], type: 'Roast', suffix: 'infused with high-energy instructions' },
  { keywords: ['zip', 'tar', 'gz', 'rar', '7z'], type: 'Pie', suffix: 'tightly packed with rich goodness' }
];

// Fallbacks
const fallbackMeals = ['Sauté', 'Curry', 'Risotto', 'Gumbo', 'Paella', 'Chowder', 'Fricassee'];
const fallbackGarnishes = ['with a dash of mystery', 'cooked to absolute perfection', 'served extra hot', 'that tastes surprisingly good'];

function determineMeal(files) {
  if (!files || files.length === 0) {
    return {
      name: "Water Soup",
      description: "A lukewarm bowl of pure air and boiled disappointment. Chef Bobby looks unimpressed.",
      color: "#add8e6",
      recipe: ["1 cup of air", "1 pinch of disappointment"]
    };
  }

  const extensionCounts = {};
  const keywords = [];
  let totalSize = 0;

  files.forEach(f => {
    totalSize += f.size;
    const ext = path.extname(f.originalname).toLowerCase().replace('.', '');
    if (ext) {
      extensionCounts[ext] = (extensionCounts[ext] || 0) + 1;
    }
    // Extract base name keywords
    const baseName = path.basename(f.originalname, path.extname(f.originalname)).toLowerCase();
    baseName.split(/[^a-zA-Z0-9]/).filter(Boolean).forEach(k => {
      if (k.length > 2) keywords.push(k);
    });
  });

  // Find most prominent extension
  let maxExt = '';
  let maxCount = 0;
  for (const [ext, count] of Object.entries(extensionCounts)) {
    if (count > maxCount) {
      maxCount = count;
      maxExt = ext;
    }
  }

  let mealType = '';
  let garnish = '';

  const matched = mealTypes.find(m => m.keywords.includes(maxExt));
  if (matched) {
    mealType = matched.type;
    garnish = matched.suffix;
  } else {
    // Pick random from fallbacks
    const h = totalSize % fallbackMeals.length;
    mealType = fallbackMeals[h];
    garnish = fallbackGarnishes[totalSize % fallbackGarnishes.length];
  }

  // Generate dynamic name
  let mainIngredient = 'Unknown File';
  if (files.length === 1) {
    mainIngredient = path.basename(files[0].originalname, path.extname(files[0].originalname));
  } else {
    // Combined name or first one
    const firstBase = path.basename(files[0].originalname, path.extname(files[0].originalname));
    mainIngredient = `${firstBase} & friends`;
  }

  // Capitalize main ingredient
  mainIngredient = mainIngredient.charAt(0).toUpperCase() + mainIngredient.slice(1);

  // Limit length
  if (mainIngredient.length > 25) {
    mainIngredient = mainIngredient.slice(0, 22) + '...';
  }

  const mealName = `${mainIngredient} ${mealType}`;

  // Custom description
  const sizeKB = (totalSize / 1024).toFixed(1);
  const description = `A grand culinary masterpiece cooked by Chef Bobby! He combined ${files.length} file(s) (${sizeKB} KB total) into a spectacular ${mealName}, ${garnish}. Warning: Ur mama might eat this too fast!`;

  // Generate ingredients
  const recipe = files.map(f => {
    const kb = (f.size / 1024).toFixed(1);
    const ext = path.extname(f.originalname).toLowerCase().replace('.', '') || 'data';
    return `- 1 pinch of ${f.originalname} (${kb} KB of pure ${ext} essence)`;
  });

  // Generate color based on size & ext
  const hash = Array.from(mealName).reduce((acc, char) => acc + char.charCodeAt(0), 0);
  const hue = hash % 360;
  const color = `hsl(${hue}, 70%, 45%)`;

  return {
    name: mealName,
    description: description,
    color: color,
    recipe: recipe
  };
}

// Simple ASCII art representation of food/beverages
const asciiArtMap = {
  'Soup': `
      (  )   (   )  )
       ) (   )  (  (
       ___________
      (___________)
       |         |
       |  SOUP   |
       \\_________/
  `,
  'Stew': `
         (  (
          )  )
       ___________
      (___________)
       \\         /
        \\_STEW__/
  `,
  'Smoothie': `
         \\_|_/
          \\ /
         [===]
         |   |
         | * |
         |___|
         (===)
  `,
  'Casserole': `
     =================
    | [][][][][][][]  |
    |  CASSEROLE      |
    |                 |
     =================
  `,
  'Stir-Fry': `
       (   )  (   )
        \\_/_/_/_/_/
        |         |
        |  FRY!   |
        \\_________/
  `,
  'Roast': `
       ,---.  ,---.
      /     \\/     \\
      \\            /
       '--.    .--'
          |    |
          |____|
  `,
  'Pie': `
         ______
       .-'      '-.
     .'  _|_|_|_   '.
     |  |_|_|_|_|   |
     '.   PIE     .'
       '-.______.-'
  `,
  'Default': `
       (\\___/)
       (='.'=)
       (")_(")
       Chef Bobby's
       Secret Recipe!
  `
};

app.post('/cook', upload.array('files'), (req, res) => {
  try {
    const files = req.files || [];
    const meal = determineMeal(files);

    // Get matching ascii art
    const baseType = meal.name.split(' ').pop();
    const asciiArt = asciiArtMap[baseType] || asciiArtMap['Default'];

    // Create unique build directory
    const buildId = `build_${Date.now()}_${Math.floor(Math.random() * 1000)}`;
    const buildDir = path.join('/tmp', buildId);
    fs.mkdirSync(buildDir, { recursive: true });

    // Generate dynamic go code
    const goCode = `
package main

import (
	"fmt"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("    YOU ARE PLAYING CHEF BOBBY'S COOKING GAME!    ")
	fmt.Println("==================================================")
	fmt.Println("")
	fmt.Println("Chef Bobby has prepared a magnificent executable meal for you:")
	fmt.Println("🍳 MEAL: ${meal.name.replace(/"/g, '\\"')}")
	fmt.Println("")
	fmt.Println("--- INGREDIENTS USED ---")
${meal.recipe.map(r => `\tfmt.Println("${r.replace(/"/g, '\\"')}")`).join('\n')}
	fmt.Println("")
	fmt.Println("--- ANALYSIS ---")
	fmt.Println("${meal.description.replace(/"/g, '\\"')}")
	fmt.Println("")
	fmt.Println("--- MEAL PRESENTATION ---")
	fmt.Println(\`${asciiArt.replace(/`/g, '\\`').replace(/\$/g, '\\$')}\`)
	fmt.Println("")
	fmt.Println("==================================================")
	fmt.Println("    *UR MAMA EATING SFX* (CRUNCH MUNCH GULP!)      ")
	fmt.Println("==================================================")
	fmt.Println("Press ENTER to finish digesting this amazing meal...")

	// Wait for user input so the window doesn't close immediately
	var input string
	fmt.Scanln(&input)
}
`;

    const mainGoPath = path.join(buildDir, 'main.go');
    fs.writeFileSync(mainGoPath, goCode);

    // Initialize module and build executable
    // We target windows x86-64 using the local standard library
    const exeName = `${meal.name.toLowerCase().replace(/[^a-z0-9]/g, '_')}.exe`;
    const exePath = path.join(buildDir, exeName);

    // Run Go build
    execSync(`go mod init bobby_meal`, { cwd: buildDir });
    execSync(`GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o "${exePath}" main.go`, { cwd: buildDir });

    // Read the compiled exe file to return to the user
    const exeBuffer = fs.readFileSync(exePath);

    // Cleanup temp files
    fs.rmSync(buildDir, { recursive: true, force: true });
    files.forEach(f => {
      try {
        fs.unlinkSync(f.path);
      } catch (e) {}
    });

    // Send file
    res.setHeader('Content-Type', 'application/octet-stream');
    res.setHeader('Content-Disposition', `attachment; filename="${exeName}"`);
    res.setHeader('X-Meal-Name', encodeURIComponent(meal.name));
    res.setHeader('X-Meal-Description', encodeURIComponent(meal.description));
    res.setHeader('X-Meal-Color', encodeURIComponent(meal.color));
    res.send(exeBuffer);

  } catch (error) {
    console.error("Error during cooking:", error);
    res.status(500).json({ error: "Chef Bobby burnt the meal! Make sure to drop valid files." });
  }
});

app.listen(PORT, () => {
  console.log(`Chef Bobby's kitchen is open at http://localhost:${PORT}`);
});
