import { describe, test, expect, mock } from "bun:test";
import React from "react";
import { renderToString } from "react-dom/server";
import { LanguageSwitcher } from "./LanguageSwitcher";

describe("LanguageSwitcher Component", () => {
  test("should render all language buttons", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="En" onLanguageChange={mockHandler} />
    );

    expect(html).toContain("En");
    expect(html).toContain("Ru");
    expect(html).toContain("Bu");
  });

  test("should apply active class to current language", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="Ru" onLanguageChange={mockHandler} />
    );

    // Check that the HTML contains the active class
    expect(html).toContain("active");
  });

  test("should render with En as current language", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="En" onLanguageChange={mockHandler} />
    );

    expect(html).toContain("En");
    expect(html).toContain("lang-button");
  });

  test("should render with Ru as current language", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="Ru" onLanguageChange={mockHandler} />
    );

    expect(html).toContain("Ru");
  });

  test("should render with Bu as current language", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="Bu" onLanguageChange={mockHandler} />
    );

    expect(html).toContain("Bu");
  });

  test("should apply correct CSS classes", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="En" onLanguageChange={mockHandler} />
    );

    expect(html).toContain("language-switcher");
    expect(html).toContain("lang-button");
  });

  test("should render three buttons", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="En" onLanguageChange={mockHandler} />
    );

    // Count button occurrences
    const buttonCount = (html.match(/lang-button/g) || []).length;
    expect(buttonCount).toBe(3);
  });

  test("should handle invalid current language gracefully", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="Invalid" onLanguageChange={mockHandler} />
    );

    // Should still render all buttons
    expect(html).toContain("En");
    expect(html).toContain("Ru");
    expect(html).toContain("Bu");
  });

  test("should render with empty current language", () => {
    const mockHandler = mock(() => {});
    const html = renderToString(
      <LanguageSwitcher currentLang="" onLanguageChange={mockHandler} />
    );

    expect(html).toContain("En");
    expect(html).toContain("Ru");
    expect(html).toContain("Bu");
  });
});
