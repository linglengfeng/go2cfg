-module(cfg_test_svr1).

-export[get/1, get_id/1, get_name/1, get_age/1, get_sex/1, get_items/1].

-include("../include/cfg.hrl").

get(1) ->
	#cfg_test_svr1{id=1,name="name2",age=10,sex=1,items=[1,2,3]};
get(2) ->
	#cfg_test_svr1{id=2,name="name2",age=10,sex=2,items=[1,2,4]};
get(_) ->
	undefined.

get_id(Val=#cfg_test_svr1{}) ->
	Val#cfg_test_svr1.id;
get_id(_) ->
	undefined.

get_name(Val=#cfg_test_svr1{}) ->
	Val#cfg_test_svr1.name;
get_name(_) ->
	undefined.

get_age(Val=#cfg_test_svr1{}) ->
	Val#cfg_test_svr1.age;
get_age(_) ->
	undefined.

get_sex(Val=#cfg_test_svr1{}) ->
	Val#cfg_test_svr1.sex;
get_sex(_) ->
	undefined.

get_items(Val=#cfg_test_svr1{}) ->
	Val#cfg_test_svr1.items;
get_items(_) ->
	undefined.


