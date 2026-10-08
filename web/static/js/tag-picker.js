class ResourceTagPicker extends HTMLElement {
  connectedCallback() {
    this.listeners?.abort();
    this.listeners = new AbortController();
    const { signal } = this.listeners;
    const details = this.querySelector("details");
    const search = this.querySelector("[data-tag-search]");
    const options = Array.from(this.querySelectorAll("[data-tag-option]"));
    const selectedTags = this.querySelector("[data-tag-selected]");
    const chips = new Map(
      Array.from(selectedTags.children, (chip) => [
        chip.getAttribute("tag-id"),
        chip,
      ]),
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
      for (const option of options) {
        const input = option.querySelector("input");
        if (input.checked && !chips.has(input.value)) {
          const tag = option.querySelector("resource-tag");
          const chip = document.createElement("resource-tag");
          chip.tag = {
            id: input.value,
            name: tag.getAttribute("name"),
            color: tag.getAttribute("color"),
          };
          chip.setAttribute("variant", "removable");
          chips.set(input.value, chip);
          selectedTags.append(chip);
        } else if (!input.checked && chips.has(input.value)) {
          chips.get(input.value).remove();
          chips.delete(input.value);
        }
      }
      selectedTags.hidden = chips.size === 0;
      this.querySelector("[data-tag-placeholder]").hidden = chips.size > 0;
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
