функция пример1(?аргумент1, ?аргумент2) {
    поток("file:./test_file_power_mem.txt") => ?поток2;
    установить_настройки_потока(?поток2, set, "mode=by_lines;");
    открыть_поток(?поток2);
    1000000 => ?целое1;
#    10 => ?целое1;
    0 => ?целое2;
    1 => ?целое3;
    "фрейм %?\r\n" => ?формат_строки;
    "фрейм_%?\r\n" => ?формат_имя_фрейма;
    "слот_%?\r\n" => ?формат_имя_слота;
    пока(?целое1 > ?целое2) {
#        печатать("значение", ?целое1);
        форматировать(?формат_имя_фрейма, ?целое1) => ?имя_фрейма;
        фрейм (наименование.имя_фрейма) => ?фрейм;
        10 => ?целое5;
        0 => ?целое6;
        1 => ?целое7;
        пока(?целое5 > ?целое6) {
             уникальный_идентификатор() => ?уникальное;
             форматировать(?формат_имя_слота, ?уникальное) => ?имя_слота;
             уникальный_идентификатор() => ?уникальное;
             добавить_слот(?фрейм, ?имя_слота);
             добавить_значение_в_слот(?фрейм, ?имя_слота, ?уникальное);
             вычесть(?целое5, ?целое7) => ?целое5;
        };
#        печатать(">> ", ?имя_фрейма);
        форматировать(?формат_строки, ?фрейм) => ?переменная_запись;
#        печатать(">> ", ?переменная_запись);
 	записать_поток(?поток2, ?переменная_запись);
# 	записать_поток(?поток2, ?имя_фрейма);
        вычесть(?целое1, ?целое3) => ?целое1;
    };
    закрыть_поток(?поток2);
};



