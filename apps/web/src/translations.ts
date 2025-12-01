export const translations = {
  En: {
    hero: {
      title: "ZIK",
      subtitle: "AI-powered development toolkit",
      description: "Commit generator, code generation, code review, auto-commenting and more.\nCompletely free and open source tool for developers.",
      github: "GitHub",
      docs: "Documentation",
    },
    stats: {
      free: "Free & Open Source",
      models: "Latest AI Models",
      unlimited: "Unlimited Usage",
      fast: "Fast & Reliable",
    },
    features: {
      title: "Features",
      commit: {
        title: "Commit Generator",
        description: "Automatic creation of meaningful commit messages based on code changes",
      },
      code: {
        title: "Code Generation",
        description: "AI assistant for writing code using the latest GLM-4 models",
      },
      review: {
        title: "Code Review",
        description: "Automatic code analysis and review with improvement suggestions",
      },
      comments: {
        title: "Auto-commenting",
        description: "Intelligent addition of comments and documentation to code",
      },
    },
    install: {
      title: "Installation",
      alternative: "Or clone the repository and build from source",
    },
    tech: {
      title: "Technologies",
      description: "Built on the latest AI models and modern technology stack",
    },
    footer: {
      created: "Created by",
      copyright: "Open Source & Free Forever",
    },
  },
  Ru: {
    hero: {
      title: "ZIK",
      subtitle: "AI-инструмент для разработчиков",
      description: "Генератор коммитов, кода, ревью, автокомментирование и многое другое.\nАбсолютно бесплатный и open source инструмент для разработчиков.",
      github: "GitHub",
      docs: "Документация",
    },
    stats: {
      free: "Бесплатно и Open Source",
      models: "Новейшие AI модели",
      unlimited: "Без ограничений",
      fast: "Быстро и надежно",
    },
    features: {
      title: "Возможности",
      commit: {
        title: "Генератор коммитов",
        description: "Автоматическое создание осмысленных commit messages на основе изменений в коде",
      },
      code: {
        title: "Генерация кода",
        description: "AI-ассистент для написания кода с использованием новейших моделей GLM-4",
      },
      review: {
        title: "Code Review",
        description: "Автоматический анализ и ревью кода с предложениями по улучшению",
      },
      comments: {
        title: "Автокомментирование",
        description: "Интеллектуальное добавление комментариев и документации к коду",
      },
    },
    install: {
      title: "Установка",
      alternative: "Или клонируйте репозиторий и соберите из исходников",
    },
    tech: {
      title: "Технологии",
      description: "Построено на новейших AI моделях и современном стеке технологий",
    },
    footer: {
      created: "Создано",
      copyright: "Open Source и бесплатно навсегда",
    },
  },
  Bu: {
    hero: {
      title: "ZIK",
      subtitle: "AI-інструмент для распрацоўшчыкаў",
      description: "Генератар камітаў, кода, рэв'ю, аўтакаментаванне і многае іншае.\nАбсалютна бясплатны і open source інструмент для распрацоўшчыкаў.",
      github: "GitHub",
      docs: "Дакументацыя",
    },
    stats: {
      free: "Бясплатна і Open Source",
      models: "Найноўшыя AI мадэлі",
      unlimited: "Без абмежаванняў",
      fast: "Хутка і надзейна",
    },
    features: {
      title: "Магчымасці",
      commit: {
        title: "Генератар камітаў",
        description: "Аўтаматычнае стварэнне асэнсаваных commit messages на аснове змен у кодзе",
      },
      code: {
        title: "Генерацыя кода",
        description: "AI-асістэнт для напісання кода з выкарыстаннем найноўшых мадэляў GLM-4",
      },
      review: {
        title: "Code Review",
        description: "Аўтаматычны аналіз і рэв'ю кода з прапановамі па паляпшэнні",
      },
      comments: {
        title: "Аўтакаментаванне",
        description: "Інтэлектуальнае дадаванне каментарыяў і дакументацыі да кода",
      },
    },
    install: {
      title: "Усталяванне",
      alternative: "Або кланіруйце рэпазіторый і сабярыце з зыходнікаў",
    },
    tech: {
      title: "Тэхналогіі",
      description: "Пабудавана на найноўшых AI мадэлях і сучасным стэку тэхналогій",
    },
    footer: {
      created: "Створана",
      copyright: "Open Source і бясплатна назаўсёды",
    },
  },
};

export type Language = keyof typeof translations;
