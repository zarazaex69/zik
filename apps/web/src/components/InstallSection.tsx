import React, { useState } from 'react';
import { Copy, Check } from 'lucide-react';

interface InstallSectionProps {
  title: string;
  alternative: string;
}

export const InstallSection: React.FC<InstallSectionProps> = ({ title, alternative }) => {
  const [copied, setCopied] = useState(false);
  const installCommand = '-fsSL https://zik.zarazaex.xyz/install | bash';

  const handleCopy = () => {
    navigator.clipboard.writeText(installCommand);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <section className="py-20 px-6">
      <div className="container mx-auto max-w-4xl">
        <h2 className="section-title">{title}</h2>
        <div className="install-card">
          <div className="flex items-center justify-between gap-4 flex-wrap">
            <code className="install-code">
              {installCommand}
            </code>
            <button
              onClick={handleCopy}
              className="copy-button"
              aria-label="Copy install command"
            >
              {copied ? (
                <Check className="w-5 h-5" />
              ) : (
                <Copy className="w-5 h-5" />
              )}
            </button>
          </div>
        </div>
        <p className="text-center text-[#888] mt-6">
          {alternative}
        </p>
      </div>
    </section>
  );
};
