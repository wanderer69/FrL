функция пример1(?аргумент1, ?аргумент2) {
    1 => ?целое1;
    2 => ?целое2;
    "12" => ?строка1;
    1.1 => ?плавучее1;
    2.6 => ?плавучее2;
    "8.0" => ?строка2;
    1 => ?целое1_1;
    3 => ?целое2_2;
    "0" => ?строка2_1;
    "." => ?строка2_2;

    сложить(?целое1, ?целое2) => ?целое3;
    печатать(?целое3);
    умножить(?целое3, ?целое2) => ?целое4;
    печатать(?целое4);
    делить(?целое4, ?целое2) => ?целое5;
    печатать(?целое5);
    вычесть(?целое4, ?целое1) => ?целое6;
    печатать(?целое6);
    из_строки(?целое6, ?строка1) => ?целое7;
    печатать(?целое7);
    сложить(?плавучее1, ?плавучее2) => ?плавучее3;
    печатать(?плавучее3);
    умножить(?плавучее3, ?плавучее2) => ?плавучее4;
    печатать(?плавучее4);
    делить(?плавучее4, ?плавучее2) => ?плавучее5;
    печатать(?плавучее5);
    вычесть(?плавучее4, ?плавучее1) => ?плавучее6;
    печатать(?плавучее6);
    из_строки(?плавучее6, ?строка2) => ?плавучее7;
    печатать(?плавучее7);
    склеить(?строка1, ?строка2) => ?строка3;
    печатать(?строка3);

    срез(?строка3, ?целое1_1, ?целое2_2) => ?строка4;
    печатать(?строка4);

    обрезать(?строка2, ?строка2_1) => ?строка5;
    печатать(?строка5);

    отрезать(?строка3, ?строка2_2) => ?строка6;
    печатать(?строка6);

    из_числа(?строка3, ?целое7) => ?строка7;
    печатать(?строка7);

    фрейм (наименование."фрейм2", сущность.объект, элемент-класс.фрейм3) => ?фрейм1;  
    для каждого элемента (?фрейм1) => ?переменная3 {
        печатать(фрейм, "значение", ?переменная3);
        получить_имя_слота(?переменная3) => ?имя_слота;
        печатать("имя слота:", ?имя_слота);
        получить_значение_слота(?переменная3) => ?значение_слота;
        печатать("значение слота:", ?значение_слота);
        получить_свойство_слота(?переменная3) => ?свойство_слота;
        печатать("свойство слота:", ?свойство_слота);

    };

};



