const tagColors = new Set([
  "red",
  "green",
  "yellow",
  "blue",
  "purple",
  "aqua",
  "gray",
  "orange",
]);

class ResourceTag extends HTMLElement {
  static observedAttributes = ["tag-id", "name", "color", "variant"];

  connectedCallback() {
    if (!this.nameElement) {
      this.innerHTML = `
        <span class="tag-label"><span class="tag-marker">#</span><span class="tag-name"></span></span>
        <button type="button" class="tag-remove" hidden>×</button>
      `;
      this.nameElement = this.querySelector(".tag-name");
      this.removeButton = this.querySelector(".tag-remove");
      this.removeButton.addEventListener("click", (event) => {
        event.preventDefault();
        event.stopPropagation();
        this.dispatchEvent(
          new CustomEvent("remove", {
            bubbles: true,
            detail: this.tag,
          }),
        );
      });
    }
    this.render();
  }

  attributeChangedCallback() {
    if (this.nameElement) this.render();
  }

  get tag() {
    return {
      id: this.getAttribute("tag-id") || "",
      name: this.getAttribute("name") || "",
      color: this.getAttribute("color") || "gray",
    };
  }

  set tag(metadata) {
    this.setAttribute("tag-id", metadata.id ?? "");
    this.setAttribute("name", metadata.name ?? "");
    this.setAttribute("color", metadata.color ?? "gray");
  }

  render() {
    const { name, color } = this.tag;
    const themeColor = tagColors.has(color) ? color : "gray";
    this.style.setProperty("--tag-color", `var(--color-${themeColor}-1)`);
    this.style.setProperty("--tag-background", `var(--color-${themeColor}-0)`);
    this.nameElement.textContent = name;
    this.removeButton.hidden = this.getAttribute("variant") !== "removable";
  }
}

customElements.define("resource-tag", ResourceTag);
