package promises

type InvalidUserError struct{

}

func (f InvalidUserError) Error() string{
	return "The user who requested this is not allowed to perform the action."
}