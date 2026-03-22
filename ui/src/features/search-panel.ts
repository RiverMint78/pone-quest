interface SearchPanelOptions {
    maxQueryChars: number;
}

interface SearchPanelController {
    clearSearchInput: () => void;
}

export function initSearchPanel(options: SearchPanelOptions): SearchPanelController {
    const searchInput: HTMLTextAreaElement | null = document.getElementById("search-input") as HTMLTextAreaElement | null;
    const clearSearch: HTMLElement | null = document.getElementById("clear-search");
    const charCount: HTMLElement | null = document.getElementById("search-char-count");

    const toChars = (text: string | null | undefined): string[] => Array.from(text ?? "");

    const updateSearchCounter = (): void => {
        if (!searchInput) return;

        let chars: string[] = toChars(searchInput.value);
        if (chars.length > options.maxQueryChars) {
            chars = chars.slice(0, options.maxQueryChars);
            searchInput.value = chars.join("");
        }

        if (charCount) {
            const current: number = chars.length;
            charCount.textContent = `${current}/${options.maxQueryChars}`;
            charCount.classList.toggle("text-brand-accent", current >= options.maxQueryChars);
            charCount.classList.toggle("font-semibold", current >= options.maxQueryChars);
        }
    };

    const clearSearchInput = (): void => {
        if (!searchInput) return;
        searchInput.value = "";
        searchInput.focus();
        updateSearchCounter();
    };

    if (searchInput) {
        updateSearchCounter();
        searchInput.addEventListener("input", updateSearchCounter);
    }

    if (clearSearch && searchInput) {
        clearSearch.addEventListener("click", clearSearchInput);
    }

    return { clearSearchInput };
}
