import random
import mojimoji
import jaconv

def random_convert(input_str):
    result = ''
    exclude_list = ['@everyone', '@here']

    i = 0
    while i < len(input_str):
        found_exclusion = any(input_str[i:i+len(exclusion)] == exclusion for exclusion in exclude_list)
        found_user_mention = input_str[i:i+3] == '<@' and '<@!'
        
        if found_exclusion or found_user_mention:
            if found_exclusion:
                exclusion = next((exclusion for exclusion in exclude_list if input_str[i:i+len(exclusion)] == exclusion), '')
                result += exclusion
                i += len(exclusion)
            else:
                # Keep the user mention as it is
                user_mention = ''
                while i < len(input_str) and input_str[i] not in ('>', ':', ' ', '<', '@', '!', '&'):
                    user_mention += input_str[i]
                    i += 1
                result += user_mention
        else:
            char = input_str[i]
            choice = random.choice(['fullwidth', 'halfwidth', 'kana', 'keep'])
            
            if choice == 'fullwidth':
                char = mojimoji.han_to_zen(char)
            elif choice == 'halfwidth':
                char = mojimoji.zen_to_han(char)
            elif choice == 'kana':
                char = jaconv.z2h(jaconv.hira2hkata(char), digit=True, ascii=True)
            
            # Convert numeric characters to halfwidth
            char = mojimoji.zen_to_han(char, digit=True, ascii=False)
            
            result += char
            i += 1

    return result
