document.addEventListener('DOMContentLoaded', () => {
    const $ = id => document.getElementById(id);
    const maxQueryChars = 500;

    const closeEpisodeViewer = () => {
        const viewer = $('episode-viewer');
        if (!viewer) return;
        viewer.innerHTML = '';
        document.body.style.overflow = '';
    };


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


    const toggleBtn = $('toggle-search'), searchContent = $('search-content'), sidebarBrand = $('sidebar-brand'), sidebarMeta = $('sidebar-meta');
    let isMobileCollapsed = false;

    const setMobileCollapsed = (collapsed) => {
        isMobileCollapsed = collapsed;

        const targets = [searchContent, sidebarBrand, sidebarMeta];
        targets.forEach((el) => {
            if (!el) return;
            el.classList.toggle('hidden', collapsed);
            el.style.display = collapsed ? 'none' : '';
        });
    };

    if (toggleBtn && searchContent) {
        toggleBtn.addEventListener('click', () => {
            setMobileCollapsed(!isMobileCollapsed);
        });
    }


    const slider = $('topk-slider'),
        realInput = $('topk-real'),
        progress = $('slider-progress'),
        display = $('topk-display');

    const steps = [10, 25, 50, 100];

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

    const searchInput = $('search-input');
    const clearSearch = $('clear-search');
    const charCount = $('search-char-count');

    const toChars = (text) => Array.from(text ?? '');

    const fallbackCopyText = (text) => {
        const temp = document.createElement('textarea');
        temp.value = text;
        temp.style.position = 'fixed';
        temp.style.top = '0';
        temp.style.left = '-9999px';
        temp.style.opacity = '0';
        document.body.appendChild(temp);
        temp.focus();
        temp.select();
        temp.setSelectionRange(0, temp.value.length);

        const ok = document.execCommand('copy');
        document.body.removeChild(temp);
        if (!ok) throw new Error('copy failed');
    };

    const copyText = async (text) => {
        if (navigator.clipboard?.writeText) {
            try {
                await navigator.clipboard.writeText(text);
                return;
            } catch (_) {
                // Fallback for mobile browsers where clipboard API is blocked.
            }
        }

        fallbackCopyText(text);
    };

    const flashCopySuccess = (button) => {
        const copyIcon = button.querySelector('.copy-icon');
        const checkIcon = button.querySelector('.check-icon');
        if (!copyIcon || !checkIcon) return;

        copyIcon.classList.add('hidden');
        checkIcon.classList.remove('hidden');

        setTimeout(() => {
            copyIcon.classList.remove('hidden');
            checkIcon.classList.add('hidden');
        }, 1200);
    };

    const updateSearchCounter = () => {
        if (!searchInput) return;

        const chars = toChars(searchInput.value);
        if (chars.length > maxQueryChars) {
            searchInput.value = chars.slice(0, maxQueryChars).join('');
        }

        if (charCount) {
            const current = toChars(searchInput.value).length;
            charCount.textContent = `${current}/${maxQueryChars}`;
            charCount.classList.toggle('text-accent', current >= maxQueryChars);
            charCount.classList.toggle('font-semibold', current >= maxQueryChars);
        }
    };

    if (searchInput) {
        updateSearchCounter();
        searchInput.addEventListener('input', updateSearchCounter);
    }

    if (clearSearch && searchInput) {
        clearSearch.addEventListener('click', () => {
            searchInput.value = '';
            searchInput.focus();
            updateSearchCounter();
        });
    }

    document.addEventListener('keydown', e => {
        if (e.key === 'Escape') {
            if ($('episode-viewer')?.children.length) {
                closeEpisodeViewer();
                return;
            }
            if (document.activeElement.id === 'search-input') {
                document.activeElement.value = '';
                updateSearchCounter();
            }
        }
    });

    let lastCopyAt = 0;

    const handleGlobalTap = (e) => {
        const copyButton = e.target.closest('[data-copy-line="true"]');
        if (copyButton) {
            const now = Date.now();
            if (now - lastCopyAt < 350) return;
            lastCopyAt = now;

            e.preventDefault();

            const lineText = copyButton.dataset.copyText ?? '';
            const character = copyButton.dataset.copyCharacter ?? 'Unknown';
            const seasonCode = copyButton.dataset.copySeasoncode ?? 'unknown';
            const payload = `${lineText}\n  —— ${character}, ${seasonCode}`;

            copyText(payload)
                .then(() => flashCopySuccess(copyButton))
                .catch(() => {
                    window.prompt('复制失败，请手动复制以下内容：', payload);
                });
            return;
        }

        if (e.target.closest('[data-episode-close="true"]')) {
            closeEpisodeViewer();
        }
    };

    document.addEventListener('click', handleGlobalTap);
    document.addEventListener('touchend', handleGlobalTap, { passive: false });

    document.body.addEventListener('htmx:afterSwap', (e) => {
        if (e.target.id !== 'episode-viewer') return;

        document.body.style.overflow = 'hidden';

        requestAnimationFrame(() => {
            const body = e.target.querySelector('.episode-modal__body');
            const highlighted = e.target.querySelector('#episode-highlight');
            if (highlighted && body) {
                const offsetTop = highlighted.offsetTop - body.offsetTop;
                const containerHeight = body.clientHeight;
                const elementHeight = highlighted.clientHeight;
                body.scrollTop = offsetTop - (containerHeight - elementHeight) / 2;
            }
        });
    });
});