import React from 'react';

interface LanguageSwitcherProps {
  currentLang: string;
  onLanguageChange: (lang: string) => void;
}

export const LanguageSwitcher: React.FC<LanguageSwitcherProps> = ({ currentLang, onLanguageChange }) => {
  const languages = ['En', 'Ru', 'Bu'];

  return (
    <div className="language-switcher">
      {languages.map((lang) => (
        <button
          key={lang}
          onClick={() => onLanguageChange(lang)}
          className={`lang-button ${currentLang === lang ? 'active' : ''}`}
        >
          {lang}
        </button>
      ))}
    </div>
  );
};
