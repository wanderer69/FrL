функция пример1(?аргумент1, ?аргумент2) {
    поток("file:../../data/scripts/lang/test_file.txt") => ?поток1;
    установить_настройки_потока(?поток1, set, "mode=full;");
    открыть_поток(?поток1);
    читать_поток(?поток1) => (?количество, ?данные);
    закрыть_поток(?поток1);
    печатать("количество", ?количество);
    печатать("значение", ?данные);
    поток("file:./test_file_new.txt") => ?поток2;
    установить_настройки_потока(?поток2, set, "mode=by_lines;");
    открыть_поток(?поток2);
    10 => ?целое1;
    0 => ?целое2;
    1 => ?целое3;
    пока(?целое1 > ?целое2) {
        печатать("значение", ?целое1);
        "строка %?\r\n" => ?строка_формата;
        форматировать(?строка_формата, ?целое1) => ?переменная7;
 	записать_поток(?поток2, ?целое1);
        вычесть(?целое1, ?целое3) => ?целое1;
    };
    закрыть_поток(?поток2);
};



