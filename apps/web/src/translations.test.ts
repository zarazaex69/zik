import { describe, test, expect } from "bun:test";
import { translations, type Language } from "./translations";

describe("Translations", () => {
  const languages: Language[] = ["En", "Ru", "Bu"];

  test("should have all required languages", () => {
    expect(Object.keys(translations)).toEqual(languages);
  });

  test("should have consistent structure across all languages", () => {
    const enKeys = Object.keys(translations.En);
    
    languages.forEach((lang) => {
      const langKeys = Object.keys(translations[lang]);
      expect(langKeys).toEqual(enKeys);
    });
  });

  test("should have hero section in all languages", () => {
    languages.forEach((lang) => {
      expect(translations[lang].hero).toBeDefined();
      expect(translations[lang].hero.title).toBeDefined();
      expect(translations[lang].hero.subtitle).toBeDefined();
      expect(translations[lang].hero.description).toBeDefined();
      expect(translations[lang].hero.github).toBeDefined();
      expect(translations[lang].hero.docs).toBeDefined();
    });
  });

  test("should have stats section in all languages", () => {
    languages.forEach((lang) => {
      expect(translations[lang].stats).toBeDefined();
      expect(translations[lang].stats.free).toBeDefined();
      expect(translations[lang].stats.models).toBeDefined();
      expect(translations[lang].stats.unlimited).toBeDefined();
      expect(translations[lang].stats.fast).toBeDefined();
    });
  });

  test("should have features section in all languages", () => {
    languages.forEach((lang) => {
      expect(translations[lang].features).toBeDefined();
      expect(translations[lang].features.title).toBeDefined();
      expect(translations[lang].features.commit).toBeDefined();
      expect(translations[lang].features.code).toBeDefined();
      expect(translations[lang].features.review).toBeDefined();
      expect(translations[lang].features.comments).toBeDefined();
    });
  });

  test("should have install section in all languages", () => {
    languages.forEach((lang) => {
      expect(translations[lang].install).toBeDefined();
      expect(translations[lang].install.title).toBeDefined();
      expect(translations[lang].install.alternative).toBeDefined();
    });
  });

  test("should have tech section in all languages", () => {
    languages.forEach((lang) => {
      expect(translations[lang].tech).toBeDefined();
      expect(translations[lang].tech.title).toBeDefined();
      expect(translations[lang].tech.description).toBeDefined();
    });
  });

  test("should have footer section in all languages", () => {
    languages.forEach((lang) => {
      expect(translations[lang].footer).toBeDefined();
      expect(translations[lang].footer.created).toBeDefined();
      expect(translations[lang].footer.copyright).toBeDefined();
    });
  });

  test("should have non-empty strings for all translations", () => {
    languages.forEach((lang) => {
      const checkObject = (obj: any, path: string = "") => {
        Object.entries(obj).forEach(([key, value]) => {
          const currentPath = path ? `${path}.${key}` : key;
          
          if (typeof value === "string") {
            expect(value.length).toBeGreaterThan(0);
          } else if (typeof value === "object" && value !== null) {
            checkObject(value, currentPath);
          }
        });
      };

      checkObject(translations[lang], lang);
    });
  });

  test("should have consistent feature structure", () => {
    const featureKeys = ["title", "description"];
    const features = ["commit", "code", "review", "comments"];

    languages.forEach((lang) => {
      features.forEach((feature) => {
        const featureObj = translations[lang].features[feature as keyof typeof translations.En.features];
        expect(Object.keys(featureObj)).toEqual(featureKeys);
      });
    });
  });

  test("English translations should be in English", () => {
    expect(translations.En.hero.title).toBe("ZIK");
    expect(translations.En.hero.github).toBe("GitHub");
    expect(translations.En.hero.docs).toBe("Documentation");
  });

  test("Russian translations should be in Russian", () => {
    expect(translations.Ru.hero.github).toBe("GitHub");
    expect(translations.Ru.hero.docs).toBe("Документация");
    expect(translations.Ru.features.title).toBe("Возможности");
  });

  test("Belarusian translations should be in Belarusian", () => {
    expect(translations.Bu.hero.github).toBe("GitHub");
    expect(translations.Bu.hero.docs).toBe("Дакументацыя");
    expect(translations.Bu.features.title).toBe("Магчымасці");
  });

  test("should have multiline description in hero", () => {
    languages.forEach((lang) => {
      expect(translations[lang].hero.description).toContain("\n");
    });
  });

  test("Language type should accept valid languages", () => {
    const validLang: Language = "En";
    expect(translations[validLang]).toBeDefined();
  });

  test("should have same number of features across languages", () => {
    const enFeaturesCount = Object.keys(translations.En.features).length;
    
    languages.forEach((lang) => {
      const langFeaturesCount = Object.keys(translations[lang].features).length;
      expect(langFeaturesCount).toBe(enFeaturesCount);
    });
  });
});
