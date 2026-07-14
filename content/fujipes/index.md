---
title: Fujipes
---

# Fujipes

X-Trans I film simulation recipes. Tap a card for the full recipe.

<div id="fujipes" class="fujipes">
  <p class="fujipes-status">Loading recipes…</p>
  <div class="fujipes-grid" hidden></div>
</div>

<dialog id="fujipes-modal" class="fujipes-modal">
  <button type="button" class="fujipes-close" aria-label="Close">&times;</button>
  <div class="fujipes-modal-body"></div>
</dialog>

<style>
.fujipes-status { color: #888; margin: 1rem 0; }
.fujipes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  margin: 1rem 0 2rem;
}
.fujipes-card {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border);
  background: #fff;
  cursor: pointer;
  text-align: left;
  font: inherit;
  color: inherit;
  padding: 0;
  overflow: hidden;
  transition: border-color 0.15s, transform 0.15s;
}
.fujipes-card:hover,
.fujipes-card:focus-visible {
  border-color: var(--header-bg);
  transform: translateY(-2px);
  outline: none;
}
.fujipes-card-img {
  aspect-ratio: 4 / 3;
  width: 100%;
  object-fit: cover;
  display: block;
  border: none;
  margin: 0;
  background: linear-gradient(145deg, #d9e2e6, #b7c5cc);
}
.fujipes-card-ph {
  aspect-ratio: 4 / 3;
  display: grid;
  place-items: center;
  background: linear-gradient(145deg, #d9e2e6, #b7c5cc);
  color: #4a6069;
  font-family: var(--heading-font-family);
  font-size: 0.75rem;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}
.fujipes-card-name {
  padding: 0.7rem 0.85rem;
  font-family: var(--heading-font-family);
  font-size: 1.05rem;
  line-height: 1.25;
}
.fujipes-modal {
  border: none;
  padding: 0;
  width: 100vw;
  max-width: 100vw;
  height: 100vh;
  max-height: 100vh;
  background: var(--bg);
  color: var(--text);
}
.fujipes-modal::backdrop { background: rgba(10, 30, 40, 0.55); }
.fujipes-close {
  position: fixed;
  top: 0.75rem;
  right: 0.75rem;
  z-index: 2;
  width: 2.5rem;
  height: 2.5rem;
  border: 1px solid var(--border);
  background: #fff;
  color: var(--text);
  font-size: 1.5rem;
  line-height: 1;
  cursor: pointer;
}
.fujipes-close:hover { border-color: var(--header-bg); }
.fujipes-modal-body {
  min-height: 100%;
  display: grid;
  grid-template-rows: minmax(40vh, 55vh) 1fr;
}
.fujipes-modal-img,
.fujipes-modal-ph {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border: none;
  margin: 0;
  background: #1a2a32;
}
.fujipes-modal-ph {
  display: grid;
  place-items: center;
  color: #c5d2d8;
  font-family: var(--heading-font-family);
  letter-spacing: 0.06em;
  text-transform: uppercase;
}
.fujipes-modal-details {
  padding: 1.25rem 1.25rem 2.5rem;
  max-width: 40rem;
  margin: 0 auto;
  width: 100%;
}
.fujipes-modal-details h2 {
  margin: 0 0 1rem;
  font-size: clamp(1.4rem, 4vw, 1.8rem);
}
.fujipes-dl {
  display: grid;
  grid-template-columns: minmax(8rem, 40%) 1fr;
  gap: 0.45rem 0.75rem;
  margin: 0;
  font-size: 1.1rem;
}
.fujipes-dl dt {
  color: #777;
  font-family: var(--heading-font-family);
  font-size: 1rem;
}
.fujipes-dl dd { margin: 0; }
.fujipes-source {
  margin-top: 1.25rem;
  font-size: 0.9rem;
  word-break: break-word;
}
@media (min-width: 700px) {
  .fujipes-modal-body {
    grid-template-rows: none;
    grid-template-columns: 1.1fr 1fr;
    min-height: 100vh;
  }
  .fujipes-modal-img,
  .fujipes-modal-ph { min-height: 100vh; }
  .fujipes-modal-details {
    padding: 3rem 2rem;
    align-self: center;
  }
}
</style>

<script>
(() => {
  const ENDPOINT = 'https://script.google.com/macros/s/AKfycbxQhmyoV3Lgkmplj2SXrEnC_nescagGvYrWqqTeCA6JmVQgmxQT_NF1OZ1OFztaux37/exec';
  const FIELDS = [
    ['ISO', 'iso'],
    ['Dynamic range', 'dynamic_range'],
    ['Film simulation', 'base_film_simulation'],
    ['Color', 'color'],
    ['Sharpness', 'sharpness'],
    ['Highlight', 'highlight_tone'],
    ['Shadow', 'shadow_tone'],
    ['Noise reduction', 'noise_reduction'],
    ['White balance', 'white_balance'],
    ['WB R', 'wb_r'],
    ['WB B', 'wb_b'],
    ['Exposure comp.', 'exposure_compensation'],
  ];

  const root = document.getElementById('fujipes');
  const status = root.querySelector('.fujipes-status');
  const grid = root.querySelector('.fujipes-grid');
  const modal = document.getElementById('fujipes-modal');
  const body = modal.querySelector('.fujipes-modal-body');
  const closeBtn = modal.querySelector('.fujipes-close');
  let recipes = [];

  const esc = (s) => String(s ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');

  const isURL = (s) => /^https?:\/\//i.test(s || '');

  function mediaHTML(r, kind) {
    if (r.image) {
      const cls = kind === 'card' ? 'fujipes-card-img' : 'fujipes-modal-img';
      return `<img class="${cls}" src="${esc(r.image)}" alt="${esc(r.name)}" loading="lazy">`;
    }
    const cls = kind === 'card' ? 'fujipes-card-ph' : 'fujipes-modal-ph';
    const label = r.base_film_simulation || 'No preview';
    return `<div class="${cls}" aria-hidden="true">${esc(label)}</div>`;
  }

  function openRecipe(i) {
    const r = recipes[i];
    if (!r) return;
    const rows = FIELDS
      .filter(([, key]) => r[key] != null && r[key] !== '')
      .map(([label, key]) => `<dt>${esc(label)}</dt><dd>${esc(r[key])}</dd>`)
      .join('');
    let source = '';
    if (r.source) {
      source = isURL(r.source)
        ? `<p class="fujipes-source"><a href="${esc(r.source)}" target="_blank" rel="noopener">Source</a></p>`
        : `<p class="fujipes-source">${esc(r.source)}</p>`;
    }
    body.innerHTML = `
      ${mediaHTML(r, 'modal')}
      <div class="fujipes-modal-details">
        <h2>${esc(r.name)}</h2>
        <dl class="fujipes-dl">${rows}</dl>
        ${source}
      </div>`;
    if (!modal.open) modal.showModal();
  }

  function render() {
    grid.innerHTML = recipes.map((r, i) => `
      <button type="button" class="fujipes-card" data-i="${i}">
        ${mediaHTML(r, 'card')}
        <span class="fujipes-card-name">${esc(r.name)}</span>
      </button>`).join('');
    status.hidden = true;
    grid.hidden = false;
  }

  grid.addEventListener('click', (e) => {
    const card = e.target.closest('.fujipes-card');
    if (card) openRecipe(+card.dataset.i);
  });
  closeBtn.addEventListener('click', () => modal.close());
  modal.addEventListener('click', (e) => {
    if (e.target === modal) modal.close();
  });
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && modal.open) modal.close();
  });

  fetch(ENDPOINT)
    .then((res) => {
      if (!res.ok) throw new Error(res.statusText);
      return res.json();
    })
    .then((json) => {
      recipes = Array.isArray(json.data) ? json.data : [];
      if (!recipes.length) {
        status.textContent = 'No recipes found.';
        return;
      }
      render();
    })
    .catch((err) => {
      status.textContent = 'Could not load recipes: ' + err.message;
    });
})();
</script>
