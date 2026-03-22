interface BackToTopOptions {
    showThresholdPx: number;
    scrollBehavior: ScrollBehavior;
}

export function initBackToTopButton(options: BackToTopOptions = {
    showThresholdPx: 120,
    scrollBehavior: "auto",
}): void {
    const button: HTMLButtonElement | null = document.getElementById("back-to-top") as HTMLButtonElement | null;
    if (!button) return;

    const getScrollTop = (): number => {
        return Math.max(window.scrollY, document.documentElement.scrollTop, document.body.scrollTop);
    };

    const syncVisibility = (): void => {
        const visible: boolean = getScrollTop() > options.showThresholdPx;
        button.classList.toggle("back-to-top-btn--visible", visible);
    };

    window.addEventListener("scroll", syncVisibility, { passive: true });
    button.addEventListener("click", (): void => {
        window.scrollTo({ top: 0, behavior: options.scrollBehavior });
    });

    syncVisibility();
}
