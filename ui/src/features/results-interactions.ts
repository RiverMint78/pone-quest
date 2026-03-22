interface ResultsInteractionOptions {
    copyDebounceMs: number;
    copyFlashMs: number;
    pinnedCardClass: string;
}

export function initResultsInteractions(options: ResultsInteractionOptions): void {
    const resultsGrid: HTMLElement | null = document.getElementById("results-grid");
    if (!resultsGrid) return;

    let lastCopyAt: number = 0;

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
        }, options.copyFlashMs);
    };

    resultsGrid.addEventListener("click", (e: Event): void => {
        const target: Element | null = e.target as Element | null;
        const copyButton: Element | null = target?.closest("[data-copy-line='true']") ?? null;
        if (copyButton) {
            const now: number = Date.now();
            if (now - lastCopyAt < options.copyDebounceMs) return;
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

        const card: HTMLElement | null = target?.closest(".result-card") as HTMLElement | null;
        if (!card) return;

        const interactiveHit: Element | null = target?.closest("button, a, input, textarea, select, label, [role='button']") ?? null;
        if (interactiveHit) return;

        card.classList.toggle(options.pinnedCardClass);
    });
}
