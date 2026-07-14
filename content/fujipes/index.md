---
title: Fujipes
---

# Fujipes

X-Trans I film simulation recipes. Tap a card for the full recipe.

<div id="fujipes" class="fujipes">
  <p class="fujipes-status">Loading recipes…</p>
  <div class="fujipes-grid" hidden></div>
  <nav class="fujipes-pager" hidden aria-label="Recipe pages"></nav>
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
.fujipes-gallery {
  position: relative;
  min-height: 0;
  background: #1a2a32;
  overflow: hidden;
}
.fujipes-modal-img,
.fujipes-modal-ph {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border: none;
  margin: 0;
  background: #1a2a32;
  display: block;
}
.fujipes-modal-ph {
  display: grid;
  place-items: center;
  color: #c5d2d8;
  font-family: var(--heading-font-family);
  letter-spacing: 0.06em;
  text-transform: uppercase;
}
.fujipes-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  z-index: 1;
  width: 2.5rem;
  height: 2.5rem;
  border: 1px solid rgba(255,255,255,0.35);
  background: rgba(0,0,0,0.45);
  color: #fff;
  font-size: 1.35rem;
  line-height: 1;
  cursor: pointer;
}
.fujipes-nav:hover { background: rgba(0,0,0,0.65); }
.fujipes-nav-prev { left: 0.75rem; }
.fujipes-nav-next { right: 0.75rem; }
.fujipes-dots {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0.85rem;
  display: flex;
  justify-content: center;
  gap: 0.4rem;
  pointer-events: none;
}
.fujipes-dot {
  width: 0.45rem;
  height: 0.45rem;
  border-radius: 50%;
  background: rgba(255,255,255,0.4);
}
.fujipes-dot.is-active { background: #fff; }
.fujipes-count {
  position: absolute;
  top: 0.85rem;
  left: 0.85rem;
  z-index: 1;
  padding: 0.2rem 0.5rem;
  background: rgba(0,0,0,0.45);
  color: #fff;
  font-family: var(--heading-font-family);
  font-size: 0.8rem;
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
.fujipes-pager {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin: 0.5rem 0 2rem;
}
.fujipes-pager button {
  font: inherit;
  font-family: var(--heading-font-family);
  font-size: 0.95rem;
  padding: 0.4rem 0.75rem;
  border: 1px solid var(--border);
  background: #fff;
  color: var(--text);
  cursor: pointer;
}
.fujipes-pager button:hover:not(:disabled) { border-color: var(--header-bg); }
.fujipes-pager button:disabled { opacity: 0.4; cursor: default; }
.fujipes-pager button.is-active {
  border-color: var(--header-bg);
  background: var(--header-bg);
  color: #fff;
}
.fujipes-pager-meta {
  color: #888;
  font-family: var(--heading-font-family);
  font-size: 0.9rem;
  margin: 0 0.25rem;
}
@media (min-width: 700px) {
  .fujipes-modal-body {
    grid-template-rows: none;
    grid-template-columns: 2fr 1fr;
    min-height: 100vh;
  }
  .fujipes-gallery { min-height: 100vh; }
  .fujipes-modal-details {
    padding: 3rem 2rem;
    align-self: center;
  }
}
</style>

<script>
(() => {
  const ENDPOINT = 'https://script.google.com/macros/s/AKfycbztLAAqQQyJfWY4YRizaPfARQdmTWADZ9Ny9-LreqfImiRSBYX62u6IBe6Wz72GBC8/exec';
  const LIMIT = 9;
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
  const pager = root.querySelector('.fujipes-pager');
  const modal = document.getElementById('fujipes-modal');
  const body = modal.querySelector('.fujipes-modal-body');
  const closeBtn = modal.querySelector('.fujipes-close');
  let recipes = [];
  let page = 1;
  let totalPages = 1;
  let slide = 0;
  let slides = [];
  let clearingURL = false;

  const esc = (s) => String(s ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');

  const isURL = (s) => /^https?:\/\//i.test(s || '');
  const slugify = (name) => String(name || '')
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '_')
    .replace(/^_|_$/g, '');
  const readSlug = () => {
    const s = location.search;
    if (s.startsWith('?=')) return decodeURIComponent(s.slice(2));
    return new URLSearchParams(s).get('') || '';
  };
  const setRecipeURL = (slug, replace = false) => {
    const next = slug
      ? `${location.pathname}?=${encodeURIComponent(slug)}`
      : location.pathname;
    if (location.pathname + location.search === next) return;
    history[replace ? 'replaceState' : 'pushState']({ recipe: slug || null }, '', next);
  };
  const toImgURL = (v) => {
    if (!v) return '';
    if (isURL(v)) return v;
    return 'https://lh3.googleusercontent.com/u/0/d/' + v;
  };
  const imgs = (r) => {
    const raw = r.images;
    if (Array.isArray(raw)) return raw.map(toImgURL).filter(Boolean);
    if (typeof raw !== 'string' || !raw.trim()) return [];
    return raw.split(/[,\s]+/).map(toImgURL).filter(Boolean);
  };

  function cardMedia(r) {
    const list = imgs(r);
    if (list[0]) {
      return `<img class="fujipes-card-img" src="${esc(list[0])}" alt="${esc(r.name)}" loading="lazy">`;
    }
    const label = r.base_film_simulation || 'No preview';
    return `<div class="fujipes-card-ph" aria-hidden="true">${esc(label)}</div>`;
  }

  function showSlide(n) {
    if (!slides.length) return;
    slide = (n + slides.length) % slides.length;
    const img = body.querySelector('.fujipes-modal-img');
    const count = body.querySelector('.fujipes-count');
    const dots = body.querySelectorAll('.fujipes-dot');
    if (img) {
      img.src = slides[slide];
      img.alt = `${body.dataset.name || ''} (${slide + 1}/${slides.length})`;
    }
    if (count) count.textContent = `${slide + 1} / ${slides.length}`;
    dots.forEach((d, i) => d.classList.toggle('is-active', i === slide));
  }

  function galleryHTML(r) {
    slides = imgs(r);
    slide = 0;
    if (!slides.length) {
      const label = r.base_film_simulation || 'No preview';
      return `<div class="fujipes-gallery"><div class="fujipes-modal-ph">${esc(label)}</div></div>`;
    }
    const multi = slides.length > 1;
    const controls = multi ? `
      <button type="button" class="fujipes-nav fujipes-nav-prev" aria-label="Previous image">‹</button>
      <button type="button" class="fujipes-nav fujipes-nav-next" aria-label="Next image">›</button>
      <div class="fujipes-count">1 / ${slides.length}</div>
      <div class="fujipes-dots">${slides.map((_, i) =>
        `<span class="fujipes-dot${i === 0 ? ' is-active' : ''}"></span>`).join('')}</div>` : '';
    return `
      <div class="fujipes-gallery">
        <img class="fujipes-modal-img" src="${esc(slides[0])}" alt="${esc(r.name)}">
        ${controls}
      </div>`;
  }

  function openRecipe(r, { syncURL = true } = {}) {
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
    body.dataset.name = r.name || '';
    body.innerHTML = `
      ${galleryHTML(r)}
      <div class="fujipes-modal-details">
        <h2>${esc(r.name)}</h2>
        <dl class="fujipes-dl">${rows}</dl>
        ${source}
      </div>`;
    if (syncURL) setRecipeURL(slugify(r.name));
    if (!modal.open) modal.showModal();
  }

  function renderPager() {
    if (totalPages <= 1) {
      pager.hidden = true;
      pager.innerHTML = '';
      return;
    }
    const nums = Array.from({ length: totalPages }, (_, i) => i + 1)
      .map((n) => `<button type="button" data-page="${n}" class="${n === page ? 'is-active' : ''}" ${n === page ? 'aria-current="page"' : ''}>${n}</button>`)
      .join('');
    pager.innerHTML = `
      <button type="button" data-page="${page - 1}" ${page <= 1 ? 'disabled' : ''}>Prev</button>
      ${nums}
      <span class="fujipes-pager-meta">${page} / ${totalPages}</span>
      <button type="button" data-page="${page + 1}" ${page >= totalPages ? 'disabled' : ''}>Next</button>`;
    pager.hidden = false;
  }

  function render() {
    grid.innerHTML = recipes.map((r, i) => `
      <button type="button" class="fujipes-card" data-i="${i}">
        ${cardMedia(r)}
        <span class="fujipes-card-name">${esc(r.name)}</span>
      </button>`).join('');
    renderPager();
    status.hidden = true;
    grid.hidden = false;
  }

  function fetchPage(n) {
    return fetch(`${ENDPOINT}?page=${n}&limit=${LIMIT}`).then((res) => {
      if (!res.ok) throw new Error(res.statusText);
      return res.json();
    });
  }

  function applyPage(json, n) {
    const meta = json.metadata || {};
    page = meta.current_page || n;
    totalPages = meta.total_pages || 1;
    recipes = Array.isArray(json.data) ? json.data : [];
  }

  function loadPage(n, { scroll = true } = {}) {
    status.hidden = false;
    status.textContent = 'Loading recipes…';
    grid.hidden = true;
    pager.hidden = true;
    return fetchPage(n)
      .then((json) => {
        applyPage(json, n);
        if (!recipes.length) {
          status.textContent = 'No recipes found.';
          return;
        }
        render();
        if (scroll) root.scrollIntoView({ behavior: 'smooth', block: 'start' });
      })
      .catch((err) => {
        status.textContent = 'Could not load recipes: ' + err.message;
      });
  }

  async function openFromSlug(slug) {
    if (!slug) return;
    let hit = recipes.find((r) => slugify(r.name) === slug);
    if (hit) {
      openRecipe(hit, { syncURL: false });
      return;
    }
    for (let p = 1; p <= totalPages; p++) {
      if (p === page && recipes.length) continue;
      try {
        const json = await fetchPage(p);
        applyPage(json, p);
        render();
        hit = recipes.find((r) => slugify(r.name) === slug);
        if (hit) {
          openRecipe(hit, { syncURL: false });
          return;
        }
      } catch (_) { /* keep looking */ }
    }
  }

  grid.addEventListener('click', (e) => {
    const card = e.target.closest('.fujipes-card');
    if (card) openRecipe(recipes[+card.dataset.i]);
  });
  pager.addEventListener('click', (e) => {
    const btn = e.target.closest('button[data-page]');
    if (!btn || btn.disabled) return;
    const n = +btn.dataset.page;
    if (n >= 1 && n <= totalPages && n !== page) loadPage(n);
  });
  body.addEventListener('click', (e) => {
    if (e.target.closest('.fujipes-nav-prev')) showSlide(slide - 1);
    if (e.target.closest('.fujipes-nav-next')) showSlide(slide + 1);
  });
  closeBtn.addEventListener('click', () => modal.close());
  modal.addEventListener('click', (e) => {
    if (e.target === modal) modal.close();
  });
  modal.addEventListener('close', () => {
    if (clearingURL) return;
    if (readSlug()) {
      clearingURL = true;
      setRecipeURL('', true);
      clearingURL = false;
    }
  });
  window.addEventListener('popstate', () => {
    const slug = readSlug();
    if (slug) openFromSlug(slug);
    else if (modal.open) {
      clearingURL = true;
      modal.close();
      clearingURL = false;
    }
  });
  document.addEventListener('keydown', (e) => {
    if (!modal.open) return;
    if (e.key === 'Escape') modal.close();
    if (e.key === 'ArrowLeft') showSlide(slide - 1);
    if (e.key === 'ArrowRight') showSlide(slide + 1);
  });

  const initialSlug = readSlug();
  loadPage(1, { scroll: false }).then(() => {
    if (initialSlug) openFromSlug(initialSlug);
  });
})();
</script>
