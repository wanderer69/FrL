пакет поиск_фрейма;
фреймы {
	фрейм (наименование."новый фрейм", сущность.объект);
	фрейм (наименование."совсем новый фрейм", сущность.субъект);
};

функция пример1(?аргумент1, ?аргумент2) {
    (сущность.объект) => ?переменная1;
    печатать(?переменная1);
};

