import { describe, test, expect, mock } from "bun:test";
import React from "react";
import { renderToString } from "react-dom/server";
import { InstallSection } from "./InstallSection";

describe("InstallSection Component", () => {
  test("should render with title and alternative text", () => {
    const title = "Quick Installation";
    const alternative = "Or download manually";

    const html = renderToString(
      <InstallSection title={title} alternative={alternative} />
    );

    expect(html).toContain("Quick Installation");
    expect(html).toContain("Or download manually");
  });

  test("should display install command", () => {
    const html = renderToString(
      <InstallSection title="Install" alternative="Alternative" />
    );

    expect(html).toContain("https://zik.zarazaex.xyz/install | bash");
  });

  test("should render copy button", () => {
    const html = renderToString(
      <InstallSection title="Install" alternative="Alternative" />
    );

    expect(html).toContain("copy-button");
  });

  test("should apply correct CSS classes", () => {
    const html = renderToString(
      <InstallSection title="Install" alternative="Alternative" />
    );

    expect(html).toContain("section-title");
    expect(html).toContain("install-card");
    expect(html).toContain("install-code");
  });

  test("should handle empty title", () => {
    const html = renderToString(
      <InstallSection title="" alternative="Alternative" />
    );

    expect(html).toContain("install-card");
  });

  test("should handle empty alternative text", () => {
    const html = renderToString(
      <InstallSection title="Install" alternative="" />
    );

    expect(html).toContain("Install");
  });

  test("should render with long title", () => {
    const longTitle = "A".repeat(200);
    const html = renderToString(
      <InstallSection title={longTitle} alternative="Alt" />
    );

    expect(html).toContain(longTitle);
  });

  test("should render with special characters in title", () => {
    const html = renderToString(
      <InstallSection
        title="Install <CLI> & Tools"
        alternative="Alternative & Options"
      />
    );

    expect(html).toBeTruthy();
  });

  test("should contain container with max-width", () => {
    const html = renderToString(
      <InstallSection title="Install" alternative="Alternative" />
    );

    expect(html).toContain("container");
    expect(html).toContain("max-w-4xl");
  });

  test("should render code element", () => {
    const html = renderToString(
      <InstallSection title="Install" alternative="Alternative" />
    );

    expect(html).toContain("<code");
  });

  test("should have aria-label on copy button", () => {
    const html = renderToString(
      <InstallSection title="Install" alternative="Alternative" />
    );

    expect(html).toContain("Copy install command");
  });
});
