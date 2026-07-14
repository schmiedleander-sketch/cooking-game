// Chef Bobby's JavaScript kitchen logic with Web Audio sound synthesis

const dropZone = document.getElementById('drop-zone');
const chefBobby = document.getElementById('chef-bobby');
const speechText = document.getElementById('speech-text');
const potLabel = document.getElementById('pot-label');
const boilingSoup = document.getElementById('boiling-soup');
const bubblesContainer = document.getElementById('bubbles');
const steamContainer = document.getElementById('steam-container');
const sfxIndicator = document.getElementById('sfx-indicator');

// Status & Progress elements
const progressBar = document.getElementById('progress-bar');
const statusText = document.getElementById('status-text');

// Modal / Presentation elements
const mealModal = document.getElementById('meal-modal');
const modalMealTitle = document.getElementById('modal-meal-title');
const modalMealDesc = document.getElementById('modal-meal-desc');
const mealIcon = document.getElementById('meal-icon');
const mealGlow = document.getElementById('meal-glow');
const btnDownload = document.getElementById('btn-download');
const btnClose = document.getElementById('btn-close');

let lastCompiledBlobUrl = null;
let lastCompiledFilename = 'meal.exe';

// Sound Synthesis using Web Audio API
const AudioContext = window.AudioContext || window.webkitAudioContext;
let audioCtx = null;

function initAudio() {
  if (!audioCtx) {
    audioCtx = new AudioContext();
  }
  if (audioCtx.state === 'suspended') {
    audioCtx.resume();
  }
}

// Generate procedurally generated bubble/boiling sound
let boilingInterval = null;
function startBoilingSound() {
  initAudio();
  if (!audioCtx) return;

  boilingInterval = setInterval(() => {
    // We synthesize a short 'plop/bubble' sound
    const osc = audioCtx.createOscillator();
    const gainNode = audioCtx.createGain();

    osc.connect(gainNode);
    gainNode.connect(audioCtx.destination);

    // Random pitch representing small and large bubbles
    const startFreq = 150 + Math.random() * 250;
    const endFreq = 400 + Math.random() * 300;

    osc.frequency.setValueAtTime(startFreq, audioCtx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(endFreq, audioCtx.currentTime + 0.15);

    gainNode.gain.setValueAtTime(0.08, audioCtx.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.001, audioCtx.currentTime + 0.15);

    osc.type = 'sine';
    osc.start();
    osc.stop(audioCtx.currentTime + 0.15);
  }, 120);
}

function stopBoilingSound() {
  if (boilingInterval) {
    clearInterval(boilingInterval);
    boilingInterval = null;
  }
}

// Procedural sound synthesis of "ur mama eating sfx" (chewing / gulping)
function playEatingSfx() {
  initAudio();
  if (!audioCtx) return;

  sfxIndicator.classList.add('sfx-active');

  // Let's queue up 5-6 chew sounds followed by a giant gulp
  let delay = 0;
  for (let i = 0; i < 6; i++) {
    setTimeout(() => {
      // Synthesize a crunch/munch sound
      const bufferSize = audioCtx.sampleRate * 0.1; // 100ms
      const buffer = audioCtx.createBuffer(1, bufferSize, audioCtx.sampleRate);
      const data = buffer.getChannelData(0);

      // Fill buffer with filtered noise for chew crunchiness
      for (let s = 0; s < bufferSize; s++) {
        data[s] = (Math.random() * 2 - 1) * Math.exp(-s / (bufferSize * 0.3));
      }

      const noiseNode = audioCtx.createBufferSource();
      noiseNode.buffer = buffer;

      const gainNode = audioCtx.createGain();
      gainNode.gain.setValueAtTime(0.2, audioCtx.currentTime);
      gainNode.gain.exponentialRampToValueAtTime(0.001, audioCtx.currentTime + 0.1);

      noiseNode.connect(gainNode);
      gainNode.connect(audioCtx.destination);
      noiseNode.start();
    }, delay);
    delay += 150 + Math.random() * 100;
  }

  // Followed by gulp sound
  setTimeout(() => {
    const osc = audioCtx.createOscillator();
    const gainNode = audioCtx.createGain();

    osc.connect(gainNode);
    gainNode.connect(audioCtx.destination);

    // Gulp: downwards pitch drop
    osc.frequency.setValueAtTime(300, audioCtx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(80, audioCtx.currentTime + 0.25);

    gainNode.gain.setValueAtTime(0.25, audioCtx.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.001, audioCtx.currentTime + 0.25);

    osc.type = 'sine';
    osc.start();
    osc.stop(audioCtx.currentTime + 0.25);

    // Remove text indicator
    setTimeout(() => {
      sfxIndicator.classList.remove('sfx-active');
    }, 400);

  }, delay + 100);
}

// Dynamically generate floating visual bubbles inside pot
let visualBubbleInterval = null;
function startVisualBubbles() {
  visualBubbleInterval = setInterval(() => {
    const bubble = document.createElement('div');
    bubble.classList.add('pot-bubble');

    const size = Math.random() * 15 + 5; // 5px to 20px
    const leftPos = Math.random() * 90; // percentage

    bubble.style.width = `${size}px`;
    bubble.style.height = `${size}px`;
    bubble.style.left = `${leftPos}%`;
    bubble.style.animationDuration = `${Math.random() * 1.2 + 0.8}s`;

    bubblesContainer.appendChild(bubble);

    // Clean up after bubble pops
    setTimeout(() => {
      bubble.remove();
    }, 2000);
  }, 100);
}

function stopVisualBubbles() {
  if (visualBubbleInterval) {
    clearInterval(visualBubbleInterval);
    visualBubbleInterval = null;
  }
}

// Drag & Drop handlers
['dragenter', 'dragover'].forEach(eventName => {
  dropZone.addEventListener(eventName, e => {
    e.preventDefault();
    e.stopPropagation();
    dropZone.classList.add('pot-hover');
    chefBobby.className = 'chef-cooking';
    speechText.innerText = "Yum! Drop them in, let's cook up a storm!";
  }, false);
});

['dragleave', 'drop'].forEach(eventName => {
  dropZone.addEventListener(eventName, e => {
    e.preventDefault();
    e.stopPropagation();
    dropZone.classList.remove('pot-hover');
  }, false);
});

dropZone.addEventListener('drop', async e => {
  const dt = e.dataTransfer;
  const files = dt.files;

  if (files && files.length > 0) {
    await cookFiles(files);
  }
});

// Alternatively, let them click to upload files
dropZone.addEventListener('click', () => {
  initAudio();
  const fileInput = document.createElement('input');
  fileInput.type = 'file';
  fileInput.multiple = true;
  fileInput.onchange = async () => {
    if (fileInput.files.length > 0) {
      await cookFiles(fileInput.files);
    }
  };
  fileInput.click();
});

// Meal visual mappings based on type
const mealEmojis = {
  'soup': '🍲',
  'stew': '🥘',
  'smoothie': '🍹',
  'casserole': '🥧',
  'stir-fry': '🥢',
  'roast': '🍖',
  'pie': '🍰',
  'default': '🍱'
};

async function cookFiles(files) {
  // Reset audio & display
  initAudio();
  chefBobby.className = 'chef-cooking';
  speechText.innerText = "Let's turn these bytes into delicious bites! Boiling the ingredients...";
  potLabel.innerText = "COOKING IN PROGRESS...";
  statusText.innerText = "Stoking the furnace...";

  startBoilingSound();
  startVisualBubbles();
  steamContainer.classList.add('steam-active');
  dropZone.classList.add('pot-cooking');

  // Multi-step visual loading progress simulation to build tension!
  const progressSteps = [
    { progress: 15, text: "Adding files to the pot..." },
    { progress: 35, text: "Dicing binary blocks..." },
    { progress: 55, text: "Boiling file headers..." },
    { progress: 75, text: "Frying compiler components..." },
    { progress: 95, text: "Packaging delicious Windows executable..." }
  ];

  for (const step of progressSteps) {
    progressBar.style.width = `${step.progress}%`;
    statusText.innerText = step.text;
    await new Promise(r => setTimeout(r, 600));
  }

  // Create form data to send to server
  const formData = new FormData();
  for (let i = 0; i < files.length; i++) {
    formData.append('files', files[i]);
  }

  try {
    const response = await fetch('/cook', {
      method: 'POST',
      body: formData
    });

    if (!response.ok) {
      const errData = await response.json();
      throw new Error(errData.error || "Chef burnt the meal.");
    }

    // Capture response headers containing meal info
    const mealName = decodeURIComponent(response.headers.get('X-Meal-Name') || 'Meal');
    const mealDesc = decodeURIComponent(response.headers.get('X-Meal-Description') || 'Fabulous meal cooked by bobby.');
    const mealColor = decodeURIComponent(response.headers.get('X-Meal-Color') || '#2ecc71');

    // Get compiled Windows EXE blob
    const blob = await response.blob();
    if (lastCompiledBlobUrl) {
      URL.revokeObjectURL(lastCompiledBlobUrl);
    }
    lastCompiledBlobUrl = URL.createObjectURL(blob);
    lastCompiledFilename = `${mealName.toLowerCase().replace(/[^a-z0-9]/g, '_')}.exe`;

    // Complete progress
    progressBar.style.width = "100%";
    statusText.innerText = "Ready to serve!";

    // Sfx
    playEatingSfx();

    // Show modal
    showMealModal(mealName, mealDesc, mealColor);

  } catch (error) {
    console.error(error);
    statusText.innerText = "Error: " + error.message;
    progressBar.style.width = "0%";
    speechText.innerText = "Oops! I burnt the files... Are you sure those ingredients are edible?";
    chefBobby.className = 'chef-idle';
    stopBoilingSound();
    stopVisualBubbles();
    steamContainer.classList.remove('steam-active');
    dropZone.classList.remove('pot-cooking');
    potLabel.innerText = "DROP INGREDIENTS (FILES) HERE";
  }
}

function showMealModal(name, desc, color) {
  // Set details
  modalMealTitle.innerText = name;
  modalMealDesc.innerText = desc;
  mealGlow.style.backgroundColor = color;

  // Set appropriate emoji
  const words = name.toLowerCase().split(' ');
  const lastWord = words[words.length - 1];
  const emoji = mealEmojis[lastWord] || mealEmojis['default'];
  mealIcon.innerText = emoji;

  // Dynamic colors for buttons
  btnDownload.style.backgroundColor = color;

  // Show
  mealModal.style.display = 'flex';
  chefBobby.className = 'chef-happy';
  speechText.innerText = `Voila! Behold the ${name}! It looks perfectly compiled!`;
}

// Download action
btnDownload.onclick = () => {
  if (lastCompiledBlobUrl) {
    const a = document.createElement('a');
    a.href = lastCompiledBlobUrl;
    a.download = lastCompiledFilename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }
};

// Close and play again
btnClose.onclick = () => {
  mealModal.style.display = 'none';
  chefBobby.className = 'chef-idle';
  speechText.innerText = "What should we cook next? Bring on more files!";
  potLabel.innerText = "DROP INGREDIENTS (FILES) HERE";
  statusText.innerText = "Stove is off";
  progressBar.style.width = "0%";

  stopBoilingSound();
  stopVisualBubbles();
  steamContainer.classList.remove('steam-active');
  dropZone.classList.remove('pot-cooking');
};
