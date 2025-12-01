import React from 'react';

interface BrutalButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  href?: string;
  className?: string;
}

export const BrutalButton: React.FC<BrutalButtonProps> = ({ children, onClick, href, className = '' }) => {
  const buttonContent = (
    <button 
      className={`brutal-button ${className}`}
      onClick={onClick}
    >
      <div className="button-overlay" />
      <span className="button-content">
        {children}
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 53 58" height={58} width={53}>
          <path strokeWidth={9} stroke="currentColor" d="M44.25 36.3612L17.25 51.9497C11.5833 55.2213 4.5 51.1318 4.50001 44.5885L4.50001 13.4115C4.50001 6.86824 11.5833 2.77868 17.25 6.05033L44.25 21.6388C49.9167 24.9104 49.9167 33.0896 44.25 36.3612Z" />
        </svg>
      </span>
    </button>
  );

  if (href) {
    return (
      <a href={href} target="_blank" rel="noopener noreferrer" className="inline-block">
        {buttonContent}
      </a>
    );
  }

  return buttonContent;
};
