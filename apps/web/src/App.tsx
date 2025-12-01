import { useState } from "react";
import "./index.css";
import { BrutalButton } from "./components/BrutalButton";
import { AnimatedBackground } from "./components/AnimatedBackground";
import { LanguageSwitcher } from "./components/LanguageSwitcher";
import { translations } from "./translations";
import type { Language } from "./translations";
import { Code2, GitCommit, MessageSquare, Sparkles, Github, ExternalLink, Copy } from "lucide-react";

export function App() {
  const [language, setLanguage] = useState<Language>("En");
  const t = translations[language];

  return (
    <div className="min-h-screen bg-[#0a0a0a] text-[#dadada]">
      <AnimatedBackground />
      
      {/* Language Switcher */}
      <div className="fixed top-6 right-6 z-50">
        <LanguageSwitcher currentLang={language} onLanguageChange={(lang) => setLanguage(lang as Language)} />
      </div>
      {/* Hero Section */}
      <section className="hero-section">
        <div className="hero-overlay" />
        
        {/* Floating decorative elements */}
        <div className="hero-decoration">
          <div className="float-element" style={{ top: '20%', left: '10%', animationDelay: '0s' }}>
            <Code2 className="w-8 h-8 opacity-20" />
          </div>
          <div className="float-element" style={{ top: '60%', right: '15%', animationDelay: '2s' }}>
            <GitCommit className="w-8 h-8 opacity-20" />
          </div>
          <div className="float-element" style={{ bottom: '30%', left: '20%', animationDelay: '4s' }}>
            <Sparkles className="w-8 h-8 opacity-20" />
          </div>
        </div>
        
        <div className="container mx-auto px-6 py-20 relative z-10">
          <div className="text-center max-w-4xl mx-auto">
            <h1 className="hero-title">
              {t.hero.title}
            </h1>
            <p className="hero-subtitle">
              {t.hero.subtitle}
            </p>
            <p className="hero-description">
              {t.hero.description.split('\n').map((line, i) => (
                <span key={i}>
                  {line}
                  {i < t.hero.description.split('\n').length - 1 && <br />}
                </span>
              ))}
            </p>
            <div className="flex flex-col gap-4 items-center mt-8">
              {/* Buttons row */}
              <div className="flex gap-4 flex-wrap justify-center">
                <BrutalButton href="https://github.com/zarazaex69/zik">
                  <Github className="w-5 h-5" />
                  {t.hero.github}
                </BrutalButton>
                <BrutalButton href="https://github.com/zarazaex69/zik#readme">
                  <ExternalLink className="w-5 h-5" />
                  {t.hero.docs}
                </BrutalButton>
              </div>
              
              {/* Install command below */}
              <div className="install-simple">
                <code className="install-command">curl zik.zarazaex.xyz/install | bash</code>
                <button
                  onClick={() => {
                    navigator.clipboard.writeText('curl zik.zarazaex.xyz/install | bash');
                  }}
                  className="install-copy"
                  title="Copy"
                >
                  <Copy className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section - Simple List */}
      <section className="py-20 px-6 relative overflow-hidden">
        <div className="container mx-auto max-w-4xl">
          <h2 className="section-title mb-16">{t.features.title}</h2>
          
          <div className="features-list">
            <div className="feature-list-item">
              <div className="feature-list-icon">
                <GitCommit className="w-6 h-6" />
              </div>
              <div className="feature-list-content">
                <h3 className="feature-list-title">{t.features.commit.title}</h3>
                <p className="feature-list-desc">{t.features.commit.description}</p>
              </div>
            </div>

            <div className="feature-list-item">
              <div className="feature-list-icon">
                <Code2 className="w-6 h-6" />
              </div>
              <div className="feature-list-content">
                <h3 className="feature-list-title">{t.features.code.title}</h3>
                <p className="feature-list-desc">{t.features.code.description}</p>
              </div>
            </div>

            <div className="feature-list-item">
              <div className="feature-list-icon">
                <MessageSquare className="w-6 h-6" />
              </div>
              <div className="feature-list-content">
                <h3 className="feature-list-title">{t.features.review.title}</h3>
                <p className="feature-list-desc">{t.features.review.description}</p>
              </div>
            </div>

            <div className="feature-list-item">
              <div className="feature-list-icon">
                <Sparkles className="w-6 h-6" />
              </div>
              <div className="feature-list-content">
                <h3 className="feature-list-title">{t.features.comments.title}</h3>
                <p className="feature-list-desc">{t.features.comments.description}</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Footer - Minimal Google-style */}
      <footer className="py-8 px-6 border-t border-[#1a1a1a]">
        <div className="container mx-auto max-w-6xl">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm">
            <div className="flex items-center gap-6 text-[#666]">
              <span>© 2025 ZIK</span>
              <a href="https://zarazaex.xyz" target="_blank" rel="noopener noreferrer" className="hover:text-[#888] transition-colors">
                zarazaex
              </a>
            </div>
            <div className="flex gap-6 text-[#666]">
              <a href="https://github.com/zarazaex69" target="_blank" rel="noopener noreferrer" className="hover:text-[#888] transition-colors">
                GitHub
              </a>
              <a href="https://t.me/zarazaex" target="_blank" rel="noopener noreferrer" className="hover:text-[#888] transition-colors">
                Telegram
              </a>
              <a href="mailto:zarazaex@tuta.io" className="hover:text-[#888] transition-colors">
                Email
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default App;
