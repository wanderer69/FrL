функция пример1(?аргумент1, ?аргумент2) {
    свойство => ?переменная2;
    фрейм (наименование."объект");
    фрейм (наименование."фрейм3", сущность.объект);  
    фрейм (наименование."фрейм2", сущность.объект, элемент-класс.фрейм3);  
    фрейм (наименование."фрейм1", сущность.объект, элемент-класс.фрейм2);  
    (сущность.объект):(?переменная2.экземпляр);
#   цикл по списку фреймов
    для каждого элемента (найти(сущность.объект)) => ?переменная3 {
        печатать(фрейм, "значение", ?переменная3);
#   цикл по списку фрейму
        для каждого элемента (?переменная3) => ?переменная4 {
            печатать(слот, "значение", ?переменная4);
        };
    };
};



