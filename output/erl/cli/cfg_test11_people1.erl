-module(cfg_test11_people1).

-export[get/1, get_id/1, get_name/1, get_age/1, get_sex/1, get_items/1].

-record(cfg, {id, name, age, sex, items}).

get(1) ->
	#cfg{id=1,name="name2",age=10,sex=1,items=[1,2,3]};
get(2) ->
	#cfg{id=2,name="name2",age=10,sex=2,items=[1,2,4]};
get(_) ->
	undefined.

get_id(Val=#cfg{}) ->
	Val#cfg.id;
get_id(_) ->
	undefined.

get_name(Val=#cfg{}) ->
	Val#cfg.name;
get_name(_) ->
	undefined.

get_age(Val=#cfg{}) ->
	Val#cfg.age;
get_age(_) ->
	undefined.

get_sex(Val=#cfg{}) ->
	Val#cfg.sex;
get_sex(_) ->
	undefined.

get_items(Val=#cfg{}) ->
	Val#cfg.items;
get_items(_) ->
	undefined.


