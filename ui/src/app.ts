import htmx from "htmx.org";

(window as Window & { htmx?: unknown }).htmx = htmx;

document.addEventListener("DOMContentLoaded", (): void => {
    const $ = (id: string): HTMLElement | null => document.getElementById(id);
    const root: HTMLElement = document.documentElement;
    const maxQueryChars: number = 500;
    const copyDebounceMs: number = 350;
    const copyFlashMs: number = 1200;
    const modalGhostClickGuardMs: number = 500;
    const minSidebarWidth: number = 320;
    const maxSidebarWidth: number = 600;
    const themeStorageKey: string = "pone-theme";

    const themeToggleBtn: HTMLElement | null = $("theme-toggle");
    const themeToggleIcon: HTMLElement | null = $("theme-toggle-icon");
        let bodyPrevOverflow: string = "";
        let bodyPrevPaddingRight: string = "";

        const lockBodyScroll = (): void => {
            const scrollbarWidth: number = window.innerWidth - document.documentElement.clientWidth;
            bodyPrevOverflow = document.body.style.overflow;
            bodyPrevPaddingRight = document.body.style.paddingRight;
            document.body.style.overflow = "hidden";
            if (scrollbarWidth > 0) {
                document.body.style.paddingRight = `${scrollbarWidth}px`;
            }
        };

        const unlockBodyScroll = (): void => {
            document.body.style.overflow = bodyPrevOverflow;
            document.body.style.paddingRight = bodyPrevPaddingRight;
        };

    const getPreferredTheme = (): "light" | "dark" => {
        const saved: string | null = window.localStorage.getItem(themeStorageKey);
        if (saved === "dark" || saved === "light") return saved;
        return "light";
    };

    const updateThemeToggleUI = (theme: "light" | "dark"): void => {
        if (themeToggleIcon) themeToggleIcon.textContent = theme === "dark" ? "☀" : "🌙";
        if (themeToggleBtn) {
            const nextModeLabel: string = theme === "dark" ? "切换到日间模式" : "切换到夜间模式";
            themeToggleBtn.setAttribute("aria-label", nextModeLabel);
            themeToggleBtn.setAttribute("title", nextModeLabel);
        }
    };

    const applyTheme = (theme: "light" | "dark"): void => {
        root.setAttribute("data-theme", theme);
        window.localStorage.setItem(themeStorageKey, theme);
        updateThemeToggleUI(theme);
    };

    applyTheme(getPreferredTheme());

    if (themeToggleBtn) {
        themeToggleBtn.addEventListener("click", (): void => {
            const current: "light" | "dark" = (root.getAttribute("data-theme") === "dark") ? "dark" : "light";
            applyTheme(current === "dark" ? "light" : "dark");
        });
    }

    const closeEpisodeViewer = (): void => {
        const viewer: HTMLElement | null = $("episode-viewer");
        if (!viewer) return;
        viewer.innerHTML = "";
           unlockBodyScroll();
    };

    const sidebar: HTMLElement | null = $("sidebar");
    const handle: HTMLElement | null = $("drag-handle");
    if (sidebar && handle) {
        let isResizing: boolean = false;
        handle.onmousedown = (e: MouseEvent): void => {
            isResizing = true;
            document.body.style.cursor = "col-resize";
            e.preventDefault();
        };
        document.onmousemove = (e: MouseEvent): void => {
            if (!isResizing) return;
            sidebar.style.width = `${Math.min(Math.max(e.clientX, minSidebarWidth), maxSidebarWidth)}px`;
        };
        document.onmouseup = (): void => {
            isResizing = false;
            document.body.style.cursor = "";
        };
    }

    const toggleBtn: HTMLElement | null = $("toggle-search");
    const searchContent: HTMLElement | null = $("search-content");
    const sidebarBrand: HTMLElement | null = $("sidebar-brand");
    const sidebarMeta: HTMLElement | null = $("sidebar-meta");
    let isMobileCollapsed: boolean = false;
    const mobileSections: (HTMLElement | null)[] = [searchContent, sidebarBrand, sidebarMeta];

    const setMobileCollapsed = (collapsed: boolean): void => {
        isMobileCollapsed = collapsed;
        mobileSections.forEach((el): void => {
            if (!el) return;
            el.classList.toggle("hidden", collapsed);
            (el as HTMLElement).style.display = collapsed ? "none" : "";
        });
    };

    if (toggleBtn && searchContent) {
        toggleBtn.addEventListener("click", () => {
            setMobileCollapsed(!isMobileCollapsed);
        });
    }

    const slider: HTMLInputElement | null = $("topk-slider") as HTMLInputElement | null;
    const realInput: HTMLInputElement | null = $("topk-real") as HTMLInputElement | null;
    const progress: HTMLElement | null = $("slider-progress") as HTMLElement | null;
    const display: HTMLElement | null = $("topk-display") as HTMLElement | null;
    const steps: number[] = [10, 25, 50, 100];

    const updateSlider = (val: string): void => {
        const idx: number = Number.parseInt(val, 10);
        const actualValue: number = steps[idx];
        if (realInput) realInput.value = String(actualValue);
        if (display) display.textContent = String(actualValue);
        if (progress) progress.style.width = `${(idx / 3) * 100}%`;
    };

    if (slider) {
        updateSlider(slider.value);
        slider.oninput = (e: Event): void => updateSlider((e.target as HTMLInputElement).value);
    }

    const searchInput: HTMLTextAreaElement | null = $("search-input") as HTMLTextAreaElement | null;
    const clearSearch: HTMLElement | null = $("clear-search");
    const charCount: HTMLElement | null = $("search-char-count");

    const toChars = (text: string | null | undefined): string[] => Array.from(text ?? "");

    const copyText = async (text: string): Promise<void> => {
        if (!navigator.clipboard?.writeText) {
            throw new Error("clipboard unavailable");
        }
        await navigator.clipboard.writeText(text);
    };

    const flashCopySuccess = (button: Element): void => {
        const copyIcon: Element | null = button.querySelector(".copy-icon");
        const checkIcon: Element | null = button.querySelector(".check-icon");
        if (!copyIcon || !checkIcon) return;
        copyIcon.classList.add("hidden");
        checkIcon.classList.remove("hidden");

        setTimeout((): void => {
            copyIcon.classList.remove("hidden");
            checkIcon.classList.add("hidden");
        }, copyFlashMs);
    };

    const updateSearchCounter = (): void => {
        if (!searchInput) return;

        let chars: string[] = toChars(searchInput.value);
        if (chars.length > maxQueryChars) {
            chars = chars.slice(0, maxQueryChars);
            searchInput.value = chars.join("");
        }

        if (charCount) {
            const current: number = chars.length;
            charCount.textContent = `${current}/${maxQueryChars}`;
            charCount.classList.toggle("text-brand-accent", current >= maxQueryChars);
            charCount.classList.toggle("font-semibold", current >= maxQueryChars);
        }
    };

    if (searchInput) {
        updateSearchCounter();
        searchInput.addEventListener("input", updateSearchCounter);
    }

    if (clearSearch && searchInput) {
        clearSearch.addEventListener("click", (): void => {
            searchInput.value = "";
            searchInput.focus();
            updateSearchCounter();
        });
    }

    document.addEventListener("keydown", (e: KeyboardEvent): void => {
        if (e.key !== "Escape") return;
        if ($("episode-viewer")?.children.length) {
            closeEpisodeViewer();
            return;
        }

        const active: HTMLTextAreaElement | null = document.activeElement as HTMLTextAreaElement | null;
        if (active?.id === "search-input") {
            active.value = "";
            updateSearchCounter();
        }
    });

    let lastCopyAt: number = 0;
    let lastModalCloseAt: number = 0;

    const handleGlobalTap = (e: Event): void => {
        if (e.type === "click" && Date.now() - lastModalCloseAt < modalGhostClickGuardMs) {
            e.preventDefault();
            return;
        }

        const target: Element | null = e.target as Element | null;
        const copyButton: Element | null = target?.closest("[data-copy-line='true']") ?? null;
        if (copyButton) {
            const now: number = Date.now();
            if (now - lastCopyAt < copyDebounceMs) return;
            lastCopyAt = now;

            e.preventDefault();

            const lineText: string = copyButton.getAttribute("data-copy-text") ?? "";
            const character: string = copyButton.getAttribute("data-copy-character") ?? "Unknown";
            const seasonCode: string = copyButton.getAttribute("data-copy-seasoncode") ?? "unknown";
            const payload: string = `${lineText}\n  -- ${character}, ${seasonCode}`;

            copyText(payload)
                .then((): void => flashCopySuccess(copyButton))
                .catch((): void => {
                    // Ignore copy failure silently.
                });
            return;
        }

        if (target?.closest("[data-episode-close='true']")) {
            lastModalCloseAt = Date.now();
            e.preventDefault();
            e.stopPropagation();
            closeEpisodeViewer();
            return;
        }
    };

    document.addEventListener("click", handleGlobalTap);
    document.addEventListener("touchend", handleGlobalTap, { passive: false });

    document.body.addEventListener("htmx:afterSwap", (e: Event): void => {
        const target: HTMLElement | null = e.target as HTMLElement | null;
        if (!target || target.id !== "episode-viewer") return;

        lockBodyScroll();
        requestAnimationFrame((): void => {
            const body: HTMLElement | null = target.querySelector(".episode-modal__body") as HTMLElement | null;
            const highlighted: HTMLElement | null = target.querySelector("#episode-highlight") as HTMLElement | null;
            const progressBar: HTMLElement | null = target.querySelector("#episode-read-progress") as HTMLElement | null;
            const hitMarker: HTMLElement | null = target.querySelector("#episode-hit-marker") as HTMLElement | null;

            const updateEpisodeReadProgress = (): void => {
                if (!body || !progressBar) return;
                const maxScrollable: number = Math.max(body.scrollHeight - body.clientHeight, 0);
                const ratio: number = maxScrollable === 0 ? 1 : Math.min(Math.max(body.scrollTop / maxScrollable, 0), 1);
                progressBar.style.width = `${Math.round(ratio * 100)}%`;
            };

            const updateHitMarkerPosition = (): void => {
                if (!body || !highlighted || !hitMarker) return;
                const maxScrollable: number = Math.max(body.scrollHeight - body.clientHeight, 1);
                const lineCenterOffset: number = (highlighted.offsetTop - body.offsetTop) + highlighted.clientHeight / 2;
                const centeredScrollTop: number = lineCenterOffset - body.clientHeight / 2;
                const ratio: number = Math.min(Math.max(centeredScrollTop / maxScrollable, 0), 1);
                hitMarker.style.left = `${Math.round(ratio * 100)}%`;
            };

            if (highlighted && body) {
                const offsetTop: number = highlighted.offsetTop - body.offsetTop;
                const containerHeight: number = body.clientHeight;
                const elementHeight: number = highlighted.clientHeight;
                body.scrollTop = offsetTop - (containerHeight - elementHeight) / 2;
            }

            updateHitMarkerPosition();
            updateEpisodeReadProgress();
            body?.addEventListener("scroll", updateEpisodeReadProgress, { passive: true });
        });
    });
});
