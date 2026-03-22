export function initTheme(root: HTMLElement, themeStorageKey = "pone-theme"): void {
    const themeToggleBtn: HTMLElement | null = document.getElementById("theme-toggle");
    const themeToggleIcon: HTMLElement | null = document.getElementById("theme-toggle-icon");

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
}
