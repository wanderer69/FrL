функция пример1(?аргумент1, ?аргумент2) {
#    свойство => ?переменная2;
    фрейм (наименование."объект") => ?переменная4;
#    фрейм (наименование."фрейм3", сущность.объект);  
#    фрейм (наименование."фрейм2", сущность.объект, элемент-класс.фрейм3);  
#    фрейм (наименование."фрейм1", сущность.объект, элемент-класс.фрейм2);  
#    (сущность.объект):(?переменная2.экземпляр);
#  ?переменная2
#  найти(сущность.объект, наименование) "123456789"
    для каждого элемента (?переменная4) => ?переменная3 {
        печатать(фрейм, "значение", ?переменная3);
    };
};
