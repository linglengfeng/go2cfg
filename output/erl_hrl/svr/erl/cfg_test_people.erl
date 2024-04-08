-module(cfg_test_people).

-export[get/1, get_id/1, get_name/1, get_age/1, get_sex/1, get_items/1, get_desc/1].

-include("../include/cfg.hrl").

get(1) ->
	#cfg_test_people{id=1,name="name1",age=10,sex=1,items=[1,2,3],desc="说明1"};
get(2) ->
	#cfg_test_people{id=2,name="name2",age=11,sex=2,items=[1,2,4],desc="说明2"};
get(_) ->
	undefined.

get_id(Val=#cfg_test_people{}) ->
	Val#cfg_test_people.id;
get_id(_) ->
	undefined.

get_name(Val=#cfg_test_people{}) ->
	Val#cfg_test_people.name;
get_name(_) ->
	undefined.

get_age(Val=#cfg_test_people{}) ->
	Val#cfg_test_people.age;
get_age(_) ->
	undefined.

get_sex(Val=#cfg_test_people{}) ->
	Val#cfg_test_people.sex;
get_sex(_) ->
	undefined.

get_items(Val=#cfg_test_people{}) ->
	Val#cfg_test_people.items;
get_items(_) ->
	undefined.

get_desc(Val=#cfg_test_people{}) ->
	Val#cfg_test_people.desc;
get_desc(_) ->
	undefined.


