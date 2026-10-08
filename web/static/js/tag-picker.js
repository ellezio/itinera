class ResourceTagPicker extends HTMLElement {
  connectedCallback() {
    this.listeners?.abort();
    this.listeners = new AbortController();
    const { signal } = this.listeners;
    const details = this.querySelector("details");
    const search = this.querySelector("[data-tag-search]");
    const options = Array.from(this.querySelectorAll("[data-tag-option]"));
    const chips = Array.from(
      this.querySelectorAll("[data-tag-selected] resource-tag"),
    );

    const filter = () => {
      const query = search.value.trim().toLocaleLowerCase();
      for (const option of options) {
        const name = option.querySelector("resource-tag").getAttribute("name");
        option.hidden = !name.toLocaleLowerCase().includes(query);
      }
      this.querySelector("[data-tag-empty]").hidden = options.some(
        (option) => !option.hidden,
      );
    };

    const updateSelection = () => {
      const selected = new Set(
        options
          .map((option) => option.querySelector("input"))
          .filter((input) => input.checked)
          .map((input) => input.value),
      );
      for (const chip of chips) {
        chip.hidden = !selected.has(chip.getAttribute("tag-id"));
      }
      this.querySelector("[data-tag-selected]").hidden = selected.size === 0;
      this.querySelector("[data-tag-placeholder]").hidden = selected.size > 0;
    };

    details.addEventListener(
      "toggle",
      () => {
        if (details.open) {
          search.focus();
        } else {
          search.value = "";
          filter();
        }
      },
      { signal },
    );

    search.addEventListener("input", filter, { signal });
    this.addEventListener("change", updateSelection, { signal });
    this.addEventListener(
      "tag-remove",
      (event) => {
        const input = options
          .map((option) => option.querySelector("input"))
          .find((input) => input.value === event.detail.id);
        if (!input) return;
        input.checked = false;
        input.dispatchEvent(new Event("change", { bubbles: true }));
      },
      { signal },
    );
    this.addEventListener(
      "keydown",
      (event) => {
        if (event.key === "Escape" && details.open) {
          event.preventDefault();
          details.open = false;
        } else if (event.target === search && event.key === "Enter") {
          // Searching must not submit the resource form.
          event.preventDefault();
        }
      },
      { signal },
    );

    document.addEventListener(
      "click",
      (event) => {
        if (!this.contains(event.target)) details.open = false;
      },
      { signal },
    );
    this.addEventListener(
      "focusout",
      (event) => {
        if (event.relatedTarget && !this.contains(event.relatedTarget)) {
          details.open = false;
        }
      },
      { signal },
    );

    filter();
    updateSelection();
  }

  disconnectedCallback() {
    this.listeners.abort();
  }
}

customElements.define("resource-tag-picker", ResourceTagPicker);
