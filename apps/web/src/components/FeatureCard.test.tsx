import { describe, test, expect } from "bun:test";
import React from "react";
import { renderToString } from "react-dom/server";
import { FeatureCard } from "./FeatureCard";

describe("FeatureCard Component", () => {
  test("should render with provided props", () => {
    const icon = <span>🚀</span>;
    const title = "Fast Performance";
    const description = "Lightning fast execution";

    const html = renderToString(
      <FeatureCard icon={icon} title={title} description={description} />
    );

    expect(html).toContain("Fast Performance");
    expect(html).toContain("Lightning fast execution");
    expect(html).toContain("🚀");
  });

  test("should render with different icon types", () => {
    const svgIcon = (
      <svg width="24" height="24">
        <circle cx="12" cy="12" r="10" />
      </svg>
    );

    const html = renderToString(
      <FeatureCard
        icon={svgIcon}
        title="Test Title"
        description="Test Description"
      />
    );

    expect(html).toContain("<svg");
    expect(html).toContain("Test Title");
  });

  test("should apply correct CSS classes", () => {
    const html = renderToString(
      <FeatureCard
        icon={<span>Icon</span>}
        title="Title"
        description="Description"
      />
    );

    expect(html).toContain("feature-card");
    expect(html).toContain("feature-icon");
    expect(html).toContain("feature-title");
    expect(html).toContain("feature-description");
  });

  test("should handle empty strings", () => {
    const html = renderToString(
      <FeatureCard icon={<span></span>} title="" description="" />
    );

    expect(html).toContain("feature-card");
  });

  test("should handle long text content", () => {
    const longTitle = "A".repeat(100);
    const longDescription = "B".repeat(500);

    const html = renderToString(
      <FeatureCard
        icon={<span>Icon</span>}
        title={longTitle}
        description={longDescription}
      />
    );

    expect(html).toContain(longTitle);
    expect(html).toContain(longDescription);
  });

  test("should render multiple instances independently", () => {
    const card1 = renderToString(
      <FeatureCard
        icon={<span>1</span>}
        title="Card 1"
        description="Description 1"
      />
    );

    const card2 = renderToString(
      <FeatureCard
        icon={<span>2</span>}
        title="Card 2"
        description="Description 2"
      />
    );

    expect(card1).toContain("Card 1");
    expect(card1).not.toContain("Card 2");
    expect(card2).toContain("Card 2");
    expect(card2).not.toContain("Card 1");
  });

  test("should handle special characters in text", () => {
    const html = renderToString(
      <FeatureCard
        icon={<span>&lt;&gt;</span>}
        title="Title with <tags>"
        description="Description with & special chars"
      />
    );

    // React escapes HTML entities
    expect(html).toBeTruthy();
  });
});
