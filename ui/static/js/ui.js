document.addEventListener('DOMContentLoaded', () => {
    const $ = id => document.getElementById(id);


    const sidebar = $('sidebar'), handle = $('drag-handle');
    if (sidebar && handle) {
        let isResizing = false;
        handle.onmousedown = (e) => { isResizing = true; document.body.style.cursor = 'col-resize'; e.preventDefault(); };
        document.onmousemove = (e) => {
            if (!isResizing) return;
            sidebar.style.width = `${Math.min(Math.max(e.clientX, 320), 600)}px`;
        };
        document.onmouseup = () => { isResizing = false; document.body.style.cursor = ''; };
    }


    const toggleBtn = $('toggle-search'), searchContent = $('search-content');
    if (toggleBtn && searchContent) {
        toggleBtn.onclick = () => searchContent.classList.toggle('hidden');
    }


    const slider = $('topk-slider'),
        realInput = $('topk-real'),
        progress = $('slider-progress'),
        display = $('topk-display');

    const steps = [5, 10, 25, 50];

    const updateSlider = (val) => {
        const idx = parseInt(val);
        const actualValue = steps[idx];

        if (realInput) realInput.value = actualValue;
        if (display) display.textContent = actualValue;
        if (progress) progress.style.width = `${(idx / 3) * 100}%`;
    };

    if (slider) {
        updateSlider(slider.value);
        slider.oninput = (e) => updateSlider(e.target.value);
    }


    const lightboxHtml = `
        <div id="lightbox" class="lightbox-modal">
            <div class="lightbox-overlay"></div>
            <div class="lightbox-content">
                <button class="lightbox-close">&times;</button>
                <img class="lightbox-image">
                <div class="lightbox-meta"></div>
            </div>
        </div>`;
    document.body.insertAdjacentHTML('beforeend', lightboxHtml);

    const lightbox = $('lightbox'),
        lbImg = lightbox.querySelector('.lightbox-image'),
        lbMeta = lightbox.querySelector('.lightbox-meta');

    const closeLB = () => {
        lightbox.classList.remove('active');
        document.body.style.overflow = '';
    };

    $('results-grid')?.addEventListener('click', e => {
        const card = e.target.closest('.result-card');
        if (!card) return;
        e.preventDefault();
        lbImg.src = card.dataset.fullImage;
        lbMeta.innerHTML = card.querySelector('.score-tag').outerHTML;
        lightbox.classList.add('active');
        document.body.style.overflow = 'hidden';
    });

    lightbox.querySelector('.lightbox-overlay').onclick = closeLB;
    lightbox.querySelector('.lightbox-close').onclick = closeLB;

    document.addEventListener('keydown', e => {
        if (e.key === 'Escape') {
            if (lightbox.classList.contains('active')) closeLB();
            else if (document.activeElement.id === 'search-input') document.activeElement.value = '';
        }
    });
});