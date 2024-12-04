# Global variable for timeout duration
TIMEOUT_DURATION="5s"  # Set this to the desired default timeout duration

wait_for_key() {
    while true; do
        read -sn 3 key < /dev/tty
        if [[ "$key" == $'\e[C' ]] || [[ "$key" == $'\e[D' ]] || [[ "$key" == $'\e[B' ]]; then
            break
        fi
    done
}

# print text with color
text() {
   local text="$1"
   local color="\033[38;5;208m"
   local RESET="\033[0m"
   printf "$(echo -e "${color}${text}${RESET}")\n"
}

run_with_timeout() {
    local file="$1"
    
    # Run the command with timeout using the global TIMEOUT_DURATION variable
    if timeout "$TIMEOUT_DURATION" ./lemin_test.sh "$file"; then
        :
    else
        text "Execution of $file timed out or was interrupted."
    fi
    wait_for_key
    rm leminTest.txt lemin.txt lemi.txt
}

test_audit() {
    text 'example0'
    run_with_timeout 'audit/example00.txt'

    text 'example1'
    run_with_timeout 'audit/example01.txt'

    text 'example2'
    run_with_timeout 'audit/example02.txt'

    text 'example3'
    run_with_timeout 'audit/example03.txt'

    text 'example4'
    run_with_timeout 'audit/example04.txt'

    text 'example5'
    run_with_timeout 'audit/example05.txt'

    text 'example6'
    run_with_timeout 'audit/example06.txt'

    text 'example7'
    run_with_timeout 'audit/example07.txt'

    text 'bad example 0'
    run_with_timeout 'audit/badexample00.txt'

    text 'bad example 1'
    run_with_timeout 'audit/badexample01.txt'
}

test_lemin() {
    text 'across'
    run_with_timeout 'cstm/across.txt'

    text 'big1'
    run_with_timeout 'cstm/big_1.txt'

    text 'big2'
    run_with_timeout 'cstm/big_2.txt'

    text 'test 1'
    run_with_timeout 'cstm/test1.txt'

    text 'test 2'
    run_with_timeout 'cstm/test2.txt'

    text 'test 3'
    run_with_timeout 'cstm/test3.txt'

    text 'example5_8'
    run_with_timeout 'cstm/exmpl5_8'

    text 'large-number'
    run_with_timeout 'cstm/large-number'

    text 'pluto1'
    run_with_timeout 'cstm/pluto_1'

    text 'pluto6'
    run_with_timeout 'cstm/pluto_6'

    text 'pluto40'
    run_with_timeout 'cstm/pluto_40'

    text 'pluto400'
    run_with_timeout 'cstm/pluto_400'

    text 'pylone1'
    run_with_timeout 'cstm/pylone_1'

    text 'pylone6'
    run_with_timeout 'cstm/pylone_6'

    text 'pylone40'
    run_with_timeout 'cstm/pylone_20'

    text 'pylone400'
    run_with_timeout 'cstm/pylone_400'
}

alias audit='test_audit'
alias lemin='test_lemin'
